package gosher

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *SshClient) download(remotePath string, localPath string) (*SshResponse, error) {
	localPathInfo, err := os.Stat(localPath)
	destinationDirectory := localPath
	var useSpecifiedFilename bool
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		} else {
			//OK - create file/dir
			useSpecifiedFilename = true
		}
	} else if localPathInfo.IsDir() {
		//ok - use name of remotePath
		destinationDirectory = localPath
		useSpecifiedFilename = false
	} else {
		destinationDirectory = filepath.Dir(localPath)
		useSpecifiedFilename = true
	}
	//from-scp
	response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)

	if err != nil {
		return response, err
	}
	defer s.session.Close()
	errorChannel := make(chan error)
	go s.manageDownloads(errorChannel, destinationDirectory, useSpecifiedFilename, localPath)
	remoteOpts := "-fr"
	err = s.session.Run("/usr/bin/scp " + remoteOpts + " " + remotePath)
	return response, err
}

func (s *SshClient) manageDownloads(errorChannel chan error, destinationDirectory string,
	useSpecifiedFilename bool, localPath string) {
	inPipe, err := s.session.StdinPipe()
	if err != nil {
		errorChannel <- err
		return
	}
	defer inPipe.Close()
	outPipe, err := s.session.StdoutPipe()
	if err != nil {
		errorChannel <- err
		return
	}
	err = sendByte(inPipe, 0)
	if err != nil {
		errorChannel <- err
		return
	}
	//use a scanner for processing individual commands, but not files themselves
	scanner := bufio.NewScanner(outPipe)
	more := true
	isFirstCommand := true
	for more {
		commandArray := make([]byte, 1)
		commandLength, err := outPipe.Read(commandArray)
		if err != nil {
			if err == io.EOF {
				//no problem.
			} else {
				errorChannel <- err
			}
			return
		}
		if commandLength < 1 {
			errorChannel <- errors.New("Error reading next byte from standard input")
			return
		}
		command := commandArray[0]
		switch command {
		case 0x0:
			//continue
		case 'E':
			//E command: go back out of dir
			destinationDirectory = filepath.Dir(destinationDirectory)
			err = sendByte(inPipe, 0)
			if err != nil {
				errorChannel <- err
				return
			}
		case 0xA:
			//0xA command: end?
			err = sendByte(inPipe, 0)
			if err != nil {
				errorChannel <- err
				return
			}
			return
		default:
			scanner.Scan()
			err = scanner.Err()
			if err != nil {
				if err == io.EOF {
					// no problem.
				} else {
					errorChannel <- err
				}
				return
			}
			// first line
			fullCommand := scanner.Text()
			// remainder, split by spaces
			splitCommands := strings.SplitN(fullCommand, " ", 3)

			switch command {
			case 0x1:
				errorChannel <- errors.New(fullCommand[1:])
				return
			case 'D', 'C':
				if err = s.manageWrites(splitCommands, inPipe, command, isFirstCommand,
					outPipe, destinationDirectory, useSpecifiedFilename, localPath, commandLength); err != nil {
					errorChannel <- err
					return
				}
			default:
				return
			}
		}
		isFirstCommand = false
	}
	err = inPipe.Close()
	if err != nil {
		errorChannel <- err
		return
	}
}

func (s *SshClient) manageWrites(splitCommands []string, inPipe io.WriteCloser, command byte, isFirstCommand bool,
	outPipe io.Reader, destinationDirectory string, useSpecifiedFilename bool, localPath string, commandLength int) error {
	mode, err := strconv.ParseInt(splitCommands[0], 8, 32)
	if err != nil {
		return err
	}
	sizeUint, err := strconv.ParseUint(splitCommands[1], 10, 64)
	commandSize := int64(sizeUint)
	if err != nil {
		return err
	}
	rcvFilename := splitCommands[2]
	var filename string
	//use the specified filename from the destination (only for top-level item)
	if useSpecifiedFilename && isFirstCommand {
		filename = filepath.Base(localPath)
	} else {
		filename = rcvFilename
	}
	err = sendByte(inPipe, 0)
	if err != nil {
		return err
	}
	if command == 'C' {
		//C command - file
		if err = s.writeFileFromPipe(destinationDirectory, filename, inPipe,
			commandLength, commandSize, outPipe); err != nil {
			return err
		}
	} else {
		//D command (directory)
		thisDstFile := filepath.Join(destinationDirectory, filename)
		fileMode := os.FileMode(uint32(mode))
		err = os.MkdirAll(thisDstFile, fileMode)
		if err != nil {
			return err
		}
		destinationDirectory = thisDstFile
	}
	return nil
}

func (s *SshClient) writeFileFromPipe(destinationDirectory string, filename string,
	inPipe io.WriteCloser, commandLength int, commandSize int64, outPipe io.Reader) error {
	thisLocalPath := filepath.Join(destinationDirectory, filename)
	tot := int64(0)

	fileWriter, err := os.Create(thisLocalPath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()

	//buffered by 4096 bytes
	bufferSize := int64(4096)
	for tot < commandSize {
		if bufferSize > commandSize-tot {
			bufferSize = commandSize - tot
		}
		b := make([]byte, bufferSize)
		commandLength, err = outPipe.Read(b)
		if err != nil {
			return err
		}
		tot += int64(commandLength)
		//write to file
		_, err = fileWriter.Write(b[:commandLength])
		if err != nil {
			return err
		}
	}
	//close file writer & check error
	err = fileWriter.Close()
	if err != nil {
		return err
	}
	//get next byte from channel reader
	nextByte := make([]byte, 1)
	_, err = outPipe.Read(nextByte)
	if err != nil {
		return err
	}
	//send null-byte back
	_, err = inPipe.Write([]byte{0})
	if err != nil {
		return err
	}
	return nil
}

func sendByte(w io.Writer, val byte) error {
	_, err := w.Write([]byte{val})
	return err
}

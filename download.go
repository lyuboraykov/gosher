package gosher

import (
	"bufio"
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *sshClient) download(remotePath string, localPath string, session *ssh.Session) (*SshResponse, error) {
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
		//localPath = filepath.Join(localPath, filepath.Base(remotePath))
		destinationDirectory = localPath
		useSpecifiedFilename = false
	} else {
		destinationDirectory = filepath.Dir(localPath)
		useSpecifiedFilename = true
	}
	//from-scp
	response := NewSshResponse(s.address, session.Stdout, session.Stderr)

	if err != nil {
		return response, err
	}
	defer session.Close()
	ce := make(chan error)
	go func() {
		cw, err := session.StdinPipe()
		if err != nil {
			ce <- err
			return
		}
		defer cw.Close()
		r, err := session.StdoutPipe()
		if err != nil {
			ce <- err
			return
		}
		err = sendByte(cw, 0)
		if err != nil {
			ce <- err
			return
		}
		//defer r.Close()
		//use a scanner for processing individual commands, but not files themselves
		scanner := bufio.NewScanner(r)
		more := true
		first := true
		for more {
			cmdArr := make([]byte, 1)
			n, err := r.Read(cmdArr)
			if err != nil {
				if err == io.EOF {
					//no problem.
				} else {
					ce <- err
				}
				return
			}
			if n < 1 {
				ce <- errors.New("Error reading next byte from standard input")
				return
			}
			cmd := cmdArr[0]
			switch cmd {
			case 0x0:
				//continue
			case 'E':
				//E command: go back out of dir
				destinationDirectory = filepath.Dir(destinationDirectory)
				err = sendByte(cw, 0)
				if err != nil {
					ce <- err
					return
				}
			case 0xA:
				//0xA command: end?
				err = sendByte(cw, 0)
				if err != nil {
					ce <- err
					return
				}

				return
			default:
				scanner.Scan()
				err = scanner.Err()
				if err != nil {
					if err == io.EOF {
						//no problem.
					} else {
						ce <- err
					}
					return
				}
				//first line
				cmdFull := scanner.Text()
				//remainder, split by spaces
				parts := strings.SplitN(cmdFull, " ", 3)

				switch cmd {
				case 0x1:
					ce <- errors.New(cmdFull[1:])
					return
				case 'D', 'C':
					mode, err := strconv.ParseInt(parts[0], 8, 32)
					if err != nil {
						ce <- err
						return
					}
					sizeUint, err := strconv.ParseUint(parts[1], 10, 64)
					size := int64(sizeUint)
					if err != nil {
						ce <- err
						return
					}
					rcvFilename := parts[2]
					var filename string
					//use the specified filename from the destination (only for top-level item)
					if useSpecifiedFilename && first {
						filename = filepath.Base(localPath)
					} else {
						filename = rcvFilename
					}
					err = sendByte(cw, 0)
					if err != nil {
						ce <- err
						return
					}
					if cmd == 'C' {
						//C command - file
						thisLocalPath := filepath.Join(destinationDirectory, filename)
						tot := int64(0)

						fw, err := os.Create(thisLocalPath)
						if err != nil {
							ce <- err
							return
						}
						defer fw.Close()

						//buffered by 4096 bytes
						bufferSize := int64(4096)
						for tot < size {
							if bufferSize > size-tot {
								bufferSize = size - tot
							}
							b := make([]byte, bufferSize)
							n, err = r.Read(b)
							if err != nil {
								ce <- err
								return
							}
							tot += int64(n)
							//write to file
							_, err = fw.Write(b[:n])
							if err != nil {
								ce <- err
								return
							}
						}
						//close file writer & check error
						err = fw.Close()
						if err != nil {
							ce <- err
							return
						}
						//get next byte from channel reader
						nb := make([]byte, 1)
						_, err = r.Read(nb)
						if err != nil {
							ce <- err
							return
						}
						//TODO check value received in nb
						//send null-byte back
						_, err = cw.Write([]byte{0})
						if err != nil {
							ce <- err
							return
						}
					} else {
						//D command (directory)
						thisDstFile := filepath.Join(destinationDirectory, filename)
						fileMode := os.FileMode(uint32(mode))
						err = os.MkdirAll(thisDstFile, fileMode)
						if err != nil {
							ce <- err
							return
						}
						destinationDirectory = thisDstFile
					}
				default:
					return
				}
			}
			first = false
		}
		err = cw.Close()
		if err != nil {
			ce <- err
			return
		}
	}()
	remoteOpts := "-fr"
	err = session.Run("/usr/bin/scp " + remoteOpts + " " + remotePath)
	return response, err
}

func sendByte(w io.Writer, val byte) error {
	_, err := w.Write([]byte{val})
	return err
}

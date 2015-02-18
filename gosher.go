// Package gosher provides types and methods for
// operations on remote machines via SSH
// e.g. execution of commands, download/upload of files
package gosher

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PASSWORD_AUTH = iota
	KEY_AUTH
)

const (
	SCP_PUSH_BEGIN_FILE   = "C0644"
	SCP_PUSH_BEGIN_FOLDER = "D0755 0"
	SCP_PUSH_END_FOLDER   = "E"
	SCP_PUSH_END          = "\x00"
)

type sshClient struct {
	address             string
	clientConfiguration ssh.ClientConfig
	Port                int
}

// Initializes the SshClient.
// This client is meant for synchronous usage with a single host.
// The client uses Port 22 by default but can be changed,
// by setting the Port field.
// address - the hostname or ip of the remote machine
// user - the username for the machine
// authenticationType - the type of authentication used, can be PASSWORD_AUTH or KEY_AUTH
// authentication - this is the password or the path to the path to the key accorrding to the authenticationType
func NewSshClient(address string, user string, authenticationType int, authentication string) (*sshClient, error) {
	if authenticationType == PASSWORD_AUTH {
		return newPasswordAuthenticatedClient(address, user, authentication), nil
	}
	keyAuthenticatedClient, err := newKeyAuthenticatedClient(address, user, authentication)
	return keyAuthenticatedClient, err
}

func newPasswordAuthenticatedClient(address string, user string, password string) *sshClient {
	configuration := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client := &sshClient{
		address:             address,
		clientConfiguration: *configuration,
		Port:                22,
	}
	return client
}

func newKeyAuthenticatedClient(address string, user string, keyPath string) (*sshClient, error) {
	key, err := getKeyFromFile(keyPath)
	if err != nil {
		return nil, err
	}
	configuration := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client := &sshClient{
		address:             address,
		clientConfiguration: *configuration,
		Port:                22,
	}
	return client, err
}

func getKeyFromFile(keyPath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return key, err
	}
	return key, err
}

func (s *sshClient) newSession() (*ssh.Session, error) {
	hostAndPort := fmt.Sprintf("%s:%d", s.address, s.Port)
	client, clientErr := ssh.Dial("tcp", hostAndPort, &s.clientConfiguration)
	if clientErr != nil {
		errorMessage := "There was an error while creating a client: " +
			clientErr.Error()
		return nil, NewSshConnectionError(errorMessage)
	}
	session, sessionErr := client.NewSession()
	if sessionErr != nil {
		errorMessage := "There was an error while establishing a session: " +
			sessionErr.Error()
		return nil, NewSshConnectionError(errorMessage)
	}
	return session, nil
}

// Executes shell command on the remote machine synchronously.
// command - the shell command to be executed on the machine.
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) ExecuteCommand(command string) (*SshResponse, error) {
	session, sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	defer session.Close()
	response := NewSshResponse(s.address, session.Stdout, session.Stderr)
	if err := session.Run(command); err != nil {
		errorMessage := "There was an error while executing the command: " +
			err.Error()
		return response, NewSshConnectionError(errorMessage)
	}
	return response, nil
}

// Executes a shell script file on the remote machine.
// It is ran in the home folder of the remote user.
// scriptPath - the path to the script on the local machine
// Returns an SshResponse and an error if any has occured
func (s *sshClient) ExecuteScript(scriptPath string) (*SshResponse, error) {
	return nil, nil
}

// Executes an function on a remote text file.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machine
// alterContentsFunction - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
// Returns SshResponse and an error if any has occured.
func (s *sshClient) ExecuteOnFile(filePath string, alterContentsFunction func(fileContent string) string) (*SshResponse, error) {
	return nil, nil
}

// Downloads file from the remote machine.
// Can be used as an alternative to scp.
// remotePath - the path to the file on the remote machine
// localPath - the path on the local machine where the file will be downloaded
// isRecursive - whether we are working with a folder or with a file
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) Download(remotePath string, localPath string, isRecursive bool) (*SshResponse, error) {
	session, sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	defer session.Close()
	return s.download(remotePath, localPath, session)
}

// Uploads file to the remote machine.
// localPath - the file on the local machine to be uploaded
// remotePath - the path on the remote machine where the file will be uploaded
// isRecursive - whether we are working with a folder or with a file
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) Upload(localPath string, remotePath string, isRecursive bool) (*SshResponse, error) {
	if isRecursive {
		return s.uploadFolder(localPath, remotePath)
	}
	return s.uploadFile(localPath, remotePath)
}

func (s *sshClient) uploadFile(localPath string, remotePath string) (*SshResponse, error) {
	session, sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	defer session.Close()
	response := NewSshResponse(s.address, session.Stdout, session.Stderr)

	go func() {
		inPipe, _ := session.StdinPipe()
		defer inPipe.Close()
		writeFileInPipe(inPipe, localPath, filepath.Base(remotePath))
	}()

	if err := session.Run("/usr/bin/scp -qvrt " + filepath.Dir(remotePath)); err != nil {
		return response, NewSshConnectionError("There was an error while uploading: " + err.Error())
	}
	return response, nil
}

func (s *sshClient) uploadFolder(localPath string, remotePath string) (*SshResponse, error) {
	session, sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	defer session.Close()
	response := NewSshResponse(s.address, session.Stdout, session.Stderr)

	go func() {
		inPipe, _ := session.StdinPipe()
		defer inPipe.Close()
		fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FOLDER, filepath.Base(remotePath))
		writeDirectoryContents(inPipe, localPath)
		fmt.Fprintln(inPipe, SCP_PUSH_END_FOLDER)
	}()

	if err := session.Run("/usr/bin/scp -qvrt " + filepath.Dir(remotePath)); err != nil {
		return response, NewSshConnectionError("Error while uploading: " + err.Error())
	}
	return response, nil
}

func writeDirectoryContents(inPipe io.WriteCloser, dir string) {
	fi, _ := ioutil.ReadDir(dir)
	for _, f := range fi {
		if f.IsDir() {
			fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FOLDER, f.Name())
			writeDirectoryContents(inPipe, dir+"/"+f.Name())
			fmt.Fprintln(inPipe, SCP_PUSH_END_FOLDER)
		} else {
			writeFileInPipe(inPipe, dir+"/"+f.Name(), f.Name())
		}
	}
}

func writeFileInPipe(inPipe io.WriteCloser, src string, remoteName string) {
	fileSrc, _ := os.Open(src)
	//Get file size
	srcStat, _ := fileSrc.Stat()
	// Print the file content
	fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FILE, srcStat.Size(), remoteName)
	io.Copy(inPipe, fileSrc)
	fmt.Fprint(inPipe, SCP_PUSH_END)
}

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

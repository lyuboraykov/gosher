// Package gosher provides types and methods for
// operations on remote machines via SSH
// e.g. execution of commands, download/upload of files
package gosher

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	PASSWORD_AUTH = iota
	KEY_AUTH
)

const (
	SCP_PUSH_BEGIN_FOLDER     = "D"
	SCP_PUSH_BEGIN_END_FOLDER = " 0"
	SCP_PUSH_END_FOLDER       = "E"
	SCP_PUSH_END              = "\x00"
)

type sshClient struct {
	hostAddress         string
	clientConfiguration ssh.ClientConfig
	Port                int
}

// Initializes the SshClient.
// This client is meant for synchronous usage with a single host.
// The client uses Port 22 by default but can be changed,
// by setting the Port field.
// hostAddress - the hostname or ip of the remote machine
// user - the username for the machine
// authenticationType - the type of authentication used, can be PASSWORD_AUTH or KEY_AUTH
// authentication - this is the password or the path to the path to the key accorrding to the authenticationType
func NewSshClient(hostAddress string, user string, authenticationType int, authentication string) (*sshClient, error) {
	if authenticationType == PASSWORD_AUTH {
		return newPasswordAuthenticatedClient(hostAddress, user, authentication), nil
	}
	keyAuthenticatedClient, err := newKeyAuthenticatedClient(hostAddress, user, authentication)
	return keyAuthenticatedClient, err
}

func newPasswordAuthenticatedClient(hostAddress string, user string, password string) *sshClient {
	configuration := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client := &sshClient{
		hostAddress:         hostAddress,
		clientConfiguration: *configuration,
		Port:                22,
	}
	return client
}

func newKeyAuthenticatedClient(hostAddress string, user string, keyPath string) (*sshClient, error) {
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
		hostAddress:         hostAddress,
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
	hostAndPort := fmt.Sprintf("%s:%d", s.hostAddress, s.Port)
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
	response := new(SshResponse)
	session.Stdout = &response.StdOut
	session.Stderr = &response.StdErr
	response.HostAddress = s.hostAddress
	if err := session.Run(command); err != nil {
		errorMessage := "There was an error while executing the command: " +
			err.Error()
		return response, NewSshConnectionError(errorMessage)
	}
	return response, nil
}

// Executes a shell script file on the remote machine.
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
	if isRecursive {
		return s.downloadFolder(localPath, remotePath)
	}
	return s.downloadFile(localPath, remotePath)
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
	response := new(SshResponse)
	session.Stdout = &response.StdOut
	session.Stderr = &response.StdErr
	response.HostAddress = s.hostAddress

	go func() {
		inPipe, _ := session.StdinPipe()
		defer inPipe.Close()

		fileSrc, _ := os.Open(localPath)

		//Get file size
		srcStat, _ := fileSrc.Stat()

		fmt.Fprintln(inPipe, "C0644", srcStat.Size(), filepath.Base(remotePath))
		io.Copy(inPipe, fileSrc)
		fmt.Fprint(inPipe, "\x00")
	}()

	if err := session.Run("/usr/bin/scp -qvrt " + filepath.Dir(remotePath)); err != nil {
		return response, NewSshConnectionError("There was an error while uploading: " + err.Error())
	}
	return response, nil
}

func getPermissions(f *os.File) (perm string) {
	fileStat, _ := f.Stat()
	mod := fileStat.Mode()
	return fmt.Sprintf("%#o", uint32(mod))
}

func (s *sshClient) uploadFolder(localPath string, remotePath string) (*SshResponse, error) {
	return nil, nil
}

func (s *sshClient) downloadFile(localPath string, remotePath string) (*SshResponse, error) {
	return nil, nil
}

func (s *sshClient) downloadFolder(localPath string, remotePath string) (*SshResponse, error) {
	return nil, nil
}

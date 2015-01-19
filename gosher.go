// Package gosher provides types and methods for
// operations on remote machines via SSH
// e.g. execution of commands, download/upload of files
package gosher

import (
	"golang.org/x/crypto/ssh"
	"ioutil"
	"strconv"
)

const (
	PASSWORD_AUTH = iota
	KEY_AUTH
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
// authentication - this is the password or the path to the key accorrding to the authenticationType
func NewSshClient(hostAddress string, user string, authenticationType int, authentication string) *sshClient {
	switch authenticationType {
	case PASSWORD_AUTH:
		return newPasswordAuthenticatedClient(hostAddress, user, authentication)
	case KEY_AUTH:
		return newKeyAuthenticatedClient(hostAddress, user, authentication)
	}
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
		clientConfiguration: configuration,
		Port:                22,
	}
	return client
}

func newKeyAuthenticatedClient(hostAddress string, user string, key string) (*sshClient, error) {
	if key, err := getKeyFile(); err != nil {
		return nil, err
	}
	configuration := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client := &sshClient{
		hostAddress:         hostAddress,
		clientConfiguration: configuration,
		Port:                22,
	}
	return client
}

func getKeyFromFile(keyPath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return key, err
	}
	return key, err
}

func (s *sshClient) newSession() (*ssh.Session, error) {
	hostAndPort := s.hostAddress + strconv.Itoa(s.Port)

	client, clientErr := Dial("tcp", hostAndPort, s.clientConfiguration)
	if clientErr != nil {
		return nil, clientErr
	}
	session, sessionErr := client.NewSession()
	if err != nil {
		return nil, sessionErr
	}
	return session, nil
}

// Executes shell command on the remote machine synchronously.
// command - the shell command to be executed on the machine.
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) ExecuteCommand(command string) (SshResponse, error) {
	session := s.newSession()
	response := new(SshResponse)
	session.Stdout = response.StdOut
	session.Stderr = response.StdErr

	if err := session.Run(command); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}

// Executes a shell script file on the remote machine.
// scriptPath - the path to the script on the local machine
// Returns an SshResponse and an error if any has occured
func (s *sshClient) ExecuteScript(scriptPath string) (SshResponse, error) {}

// Executes an function on a remote text file.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machine
// fn - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
// Returns SshResponse and an error if any has occured.
func (s *sshClient) ExecuteOnFile(filePath string, fn func(fileContent string) string) (SshResponse, error) {
}

// Downloads file from the remote machine.
// Can be used as an alternative to scp.
// remotePath - the path to the file on the remote machine
// localPath - the path on the local machine where the file will be downloaded
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) DownloadFile(remotePath string, localPath string) (SshResponse, error) {}

// Uploads file to the remote machine.
// localPath - the file on the local machine to be uploaded
// remotePath - the path on the remote machine where the file will be uploaded
// Returns an SshResponse and an error if any has occured.
func (s *sshClient) UploadFile(localPath string, remotePath string) (SshResponse, error) {}

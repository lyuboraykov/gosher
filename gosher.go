// Package gosher provides types and methods for
// operations on remote machines via SSH
// e.g. execution of commands, download/upload of files
package gosher

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
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

// Port - 22 by default
// StickySession - false by default, if true
// sessions won't be closed automatically and one would have to use
// CloseSession()
type SshClient struct {
	Port                int
	StickySession       bool
	address             string
	clientConfiguration ssh.ClientConfig
	session             ssh.Session
	isSessionOpened     bool
}

// Initializes the SshClient.
// This client is meant for synchronous usage with a single host.
// address - the hostname or ip of the remote machine
// user - the username for the machine
// authenticationType - the type of authentication used, can be PASSWORD_AUTH or KEY_AUTH
// authentication - this is the password or the path to the path to the key accorrding to the authenticationType
func NewSshClient(address string, user string, authenticationType int, authentication string) (*SshClient, error) {
	if authenticationType == PASSWORD_AUTH {
		return newPasswordAuthenticatedClient(address, user, authentication), nil
	}
	keyAuthenticatedClient, err := newKeyAuthenticatedClient(address, user, authentication)
	return keyAuthenticatedClient, err
}

func newPasswordAuthenticatedClient(address string, user string, password string) *SshClient {
	configuration := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client := &SshClient{
		address:             address,
		clientConfiguration: *configuration,
		Port:                22,
		StickySession:       false,
		isSessionOpened:     false,
	}
	return client
}

func newKeyAuthenticatedClient(address string, user string, keyPath string) (*SshClient, error) {
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
	client := &SshClient{
		address:             address,
		clientConfiguration: *configuration,
		Port:                22,
		StickySession:       false,
		isSessionOpened:     false,
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

func (s *SshClient) newSession() error {
	if !s.isSessionOpened {
		hostAndPort := fmt.Sprintf("%s:%d", s.address, s.Port)
		client, clientErr := ssh.Dial("tcp", hostAndPort, &s.clientConfiguration)
		if clientErr != nil {
			errorMessage := "There was an error while creating a client: " +
				clientErr.Error()
			return NewSshConnectionError(errorMessage)
		}
		session, sessionErr := client.NewSession()
		if sessionErr != nil {
			errorMessage := "There was an error while establishing a session: " +
				sessionErr.Error()
			return NewSshConnectionError(errorMessage)
		}
		s.isSessionOpened = true
		s.session = *session
	}
	return nil
}

// Executes shell command on the remote machine synchronously.
// command - the shell command to be executed on the machine.
// Returns an SshResponse and an error if any has occured.
func (s *SshClient) Run(command string) (*SshResponse, error) {
	sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	if !s.StickySession {
		defer s.CloseSession()
	}
	response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)
	if err := s.session.Run(command); err != nil {
		errorMessage := "There was an error while executing the command: " +
			err.Error()
		return response, NewSshConnectionError(errorMessage)
	}
	return response, nil
}

// Executes a shell script file on the remote machine.
// It is copied in the tmp folder and ran in a single session.
// chmod +x is applied before running.
// scriptPath - the path to the script on the local machine
// Returns an SshResponse and an error if any has occured
func (s *SshClient) RunScript(scriptPath string) (*SshResponse, error) {
	sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	if !s.StickySession {
		defer s.CloseSession()
	}
	response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)
	remotePath := fmt.Sprintf("/tmp/%s", filepath.Base(scriptPath))
	if _, upErr := s.uploadFile(scriptPath, remotePath); upErr != nil {
      return response, upErr
   }
	executeCommand := fmt.Sprintf("chmod +x %s ; %s", remotePath, remotePath)
	if err := s.session.Run(executeCommand); err != nil {
		errorMessage := "There was an error while executing the script: " +
			err.Error()
		return response, NewSshConnectionError(errorMessage)
	}
	return response, nil
}

// Executes an function on a remote text file.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machine
// alterContentsFunction - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
// Returns SshResponse and an error if any has occured.
func (s *SshClient) RunOnFile(filePath string, alterContentsFunction func(fileContent string) string) (*SshResponse, error) {
	sessionErr := s.newSession()
   if sessionErr != nil {
      return nil, sessionErr
   }
   if !s.StickySession {
      defer s.CloseSession()
   }
   response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)
   temporaryLocalPath := fmt.Sprintf("/tmp/%s", filepath.Base(filePath))
   if _, downloadErr := s.download(filePath, temporaryLocalPath); downloadErr != nil {
      return nil, downloadErr
   }
   buf, err := ioutil.ReadFile(temporaryLocalPath)
   if err != nil {
      return nil, err
   }
   fileContent := string(buf)
   newFileContent := alterContentsFunction(fileContent)
   ioutil.WriteFile(temporaryLocalPath, []byte(newFileContent), os.ModeTemporary)
   if runErr := s.session.Run("rm -f " + filePath); runErr != nil {
      return nil, runErr
   }
   if _, upErr := s.uploadFile(temporaryLocalPath, filePath); upErr != nil {
      return nil, upErr
   }
   return response, nil
}

// Downloads file from the remote machine.
// Can be used as an alternative to scp.
// remotePath - the path to the file on the remote machine
// localPath - the path on the local machine where the file will be downloaded
// Returns an SshResponse and an error if any has occured.
func (s *SshClient) Download(remotePath string, localPath string) (*SshResponse, error) {
	sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	if !s.StickySession {
		defer s.CloseSession()
	}
	return s.download(remotePath, localPath)
}

// Uploads file to the remote machine.
// localPath - the file on the local machine to be uploaded
// remotePath - the path on the remote machine where the file will be uploaded
// isRecursive - whether we are working with a folder or with a file
// Returns an SshResponse and an error if any has occured.
func (s *SshClient) Upload(localPath string, remotePath string) (*SshResponse, error) {
	localPathInfo, err := os.Stat(localPath)
	if err != nil {
		return nil, err
	}
	sessionErr := s.newSession()
	if sessionErr != nil {
		return nil, sessionErr
	}
	if !s.StickySession {
		defer s.CloseSession()
	}
	if localPathInfo.IsDir() {
		return s.uploadFolder(localPath, remotePath)
	} else {
		return s.uploadFile(localPath, remotePath)
	}
}

// Closes the session, use only with StickySession set to true
func (s *SshClient) CloseSession() error {
	s.isSessionOpened = false
	s.session.Close()
	return nil
}

package gosher

import (
	"golang.org/x/crypto/ssh"
	"ioutil"
)

const (
	PASSWORD_AUTH = iota
	KEY_AUTH
)

type sshClient struct {
	host                string
	clientConfiguration ssh.ClientConfig
}

func NewSshClient(host string, user string, authenticationType int, authentication string) *sshClient {
	switch authenticationType {
	case PASSWORD_AUTH:
		return newPasswordAuthenticatedClient(host, user, authentication)
	case KEY_AUTH:
		return newKeyAuthenticatedClient(host, user, authentication)
	}
}

func newPasswordAuthenticatedClient(host string, user string, password string) *sshClient {
	configuration := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	client := &sshClient{
		host:                host,
		clientConfiguration: configuration,
	}
	return client
}

func newKeyAuthenticatedClient(host string, user string, key string) (*sshClient, error) {
	if key, err := getKeyFile(); err != nil {
		return _, err
	}
	configuration := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client := &sshClient{
		host:                host,
		clientConfiguration: configuration,
	}
	return client
}

func getKeyFromFile(keyPath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return _, err
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return key, err
	}
	return key, err
}

func (s *sshClient) ExecuteCommand(command string) (SshResponse, error)   {}
func (s *sshClient) ExecuteScript(scriptPath string) (SshResponse, error) {}
func (s *sshClient) ExecuteOnFile(filePath string, fn func(fileContent string) string) (SshResponse, error) {
}
func (s *sshClient) DownloadFile(remotePath string, localPath string) (SshResponse, error) {}
func (s *sshClient) UploadFile(localPath string, remotePath string) (SshResponse, error)   {}

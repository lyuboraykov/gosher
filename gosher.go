package gosher

import "golang.org/x/crypto/ssh"

const (
	PASSWORD_AUTH = iota
	KEY_AUTH
	KEY_PATH_AUTH
)

type SshClient struct {
	host                string
	clientConfiguration ssh.ClientConfig
	keepAlive           bool
}

func NewSshClient(host string, authenticationType int, authentication string, keepAlive bool) *SshClient {
	switch authenticationType {
	case PASSWORD_AUTH:
		return newPasswordAuthenticatedClient(host, authentication)
	case KEY_AUTH:
		return newKeyAuthenticatedClient(host, authentication)
	case KEY_PATH_AUTH:
		key := getKeyFromFile(authentication)
		return newKeyAuthenticatedClient(host, key)
	}
}

func newPasswordAuthenticatedClient(host string, password string) *SshClient {}

func newKeyAuthenticatedClient(host string, key string) *SshClient {}

func getKeyFromFile(keyPath string) string {}

func (s *SshClient) ExecuteCommand(command string) (SshResponse, error)   {}
func (s *SshClient) ExecuteScript(scriptPath string) (SshResponse, error) {}
func (s *SshClient) ExecuteOnFile(filePath string, fn func(fileContent string) string) (SshResponse, error) {
}
func (s *SshClient) DownloadFile(remotePath string, localPath string) (SshResponse, error) {}
func (s *SshClient) UploadFile(localPath string, remotePath string) (SshResponse, error)   {}

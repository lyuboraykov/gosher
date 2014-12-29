package gosher

const (
    PASSWORD_AUTH = iota
    KEY_AUTH
    KEY_PATH_AUTH
)

type sshConnection struct{
   host string
   authenticationType int
   authentication string
   keepAlive bool
}

type sshResponse struct {
	errorMessage string
   exitCode int
   stdOut string
   stdErr string
}

func (sr *SshResponse) Error() string {
	return sr.errorMessage
}

func NewSshConnection() *SshConnection                                                           {}
func (s *SshConnection) ExecuteCommand(command string) SshResponse                               {}
func (s *SshConnection) ExecuteScript(scriptPath string) SshResponse                             {}
func (s *SshConnection) ExecuteOnFile(filePath string, fn func(fileContents string)) SshResponse {}

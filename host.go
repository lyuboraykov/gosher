package gosher

// Host - remote machine definition type.
// Used with the multipleHostsSshClient with asynchronous execution.
// Client - the SshClient responsible for this host
// ResultChannel - the channel via which the SshResponse of the operations will be passed
// ErrorChannel - the channel via which the error of the operations will be passed
type Host struct {
	Client        *SshClient
	ResultChannel chan *SshResponse
	ErrorChannel  chan error
}

// Constructor method for the Host type
func NewHost(address string, user string, authenticationType int, authentication string,
	resultChannel chan *SshResponse, errorChannel chan error) *Host {
	host := new(Host)
	host.Client = NewSshClient(address, user, authenticationType, authentication)
	host.ResultChannel = resultChannel
	host.ErrorChannel = errorChannel
	return &host
}

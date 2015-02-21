package gosher

// Host - remote machine definition type.
// Use with the multipleHostsSshClient with asynchronous execution.
// ResultChannel - the channel via which the SshResponse of the operations will be passed.
// ErrorChannel - the channel via which the error of the operations will be passed.
type Host struct {
	Client        *SshClient
	ResultChannel chan *SshResponse
	ErrorChannel  chan error
}

// Constructor method for the Host type
func NewHost(address string, user string, authenticationType int, authentication string,
	resultChannel chan *SshResponse, errorChannel chan error) (*Host, error) {
	host := new(Host)
	client, err := NewSshClient(address, user, authenticationType, authentication)
	if err != nil {
		return nil, err
	}
	host.Client = client
	host.ResultChannel = resultChannel
	host.ErrorChannel = errorChannel
	return host, nil
}

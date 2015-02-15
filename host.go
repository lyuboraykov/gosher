package gosher

// Host - remote machine definition type.
// Used with the multipleHostsSshClient with asynchronous execution.
// hostAddress - the address of the remote machine
// authenticationType - PASSWORD_AUTH or KEY_AUTH
// authentication - either the password of the machine or the path to the key on
// the local machine, according to the authentication type
// resultChannel - the channel via which the SshResponse of the operations will be
// passed
type Host struct {
	hostAddress        string
	authenticationType int
	authentication     string
	resultChannel      chan *SshResponse
}

// Constructor method for the Host type
func NewHost(hostAddress string, authenticationType int, authentication string,
	resultChannel chan *SshResponse) *Host {
	return nil
}

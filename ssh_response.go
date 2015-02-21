package gosher

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

// Standard response returned from ssh operations
type SshResponse struct {
	Address string
	StdOut  bytes.Buffer
	StdErr  bytes.Buffer
}

func NewSshResponse(host string, session *ssh.Session) *SshResponse {
	response := new(SshResponse)
	response.Address = host
	session.Stdout = &response.StdOut
	session.Stderr = &response.StdErr
	return response
}

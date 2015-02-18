package gosher

import "bytes"
import "io"

// Standard response returned from ssh operations
// Address - the address/host of the machine the operation was executed on
// ExitCode - the exit code of the exeecuted command.
// StdOut - the standard output of the command.
// StdErr - the standard error stream of the command.
type SshResponse struct {
	Address string
	StdOut  bytes.Buffer
	StdErr  bytes.Buffer
}

func NewSshResponse(host string, sessionStdout io.Writer, sessionStderr io.Writer) *SshResponse {
	response := new(SshResponse)
	response.Address = host
	sessionStdout = &response.StdOut
	sessionStderr = &response.StdErr
	return response
}

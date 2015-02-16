package gosher

import "bytes"
import "io"

// Standard response returned from ssh operations
// HostAddress - the address/host of the machine the operation was executed on
// ExitCode - the exit code of the exeecuted command.
// StdOut - the standard output of the command.
// StdErr - the standard error stream of the command.
type SshResponse struct {
	HostAddress string
	StdOut      bytes.Buffer
	StdErr      bytes.Buffer
}

func NewSshResponse(host string, sessionStdout io.Writer, sessionStderr io.Writer) *SshResponse {
	response := new(SshResponse)
	response.HostAddress = host
	sessionStdout = &response.StdOut
	sessionStderr = &response.StdErr
	return response
}

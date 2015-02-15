package gosher

import "bytes"

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

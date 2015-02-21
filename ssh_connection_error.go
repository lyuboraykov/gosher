package gosher

// Standard error returned on all ssh operations
// This means there was an error with the connection or the command
// returned an error code different from 0.
type SshConnectionError struct {
	errorMessage string
}

// Returns the error message of the SshConnectionError
func (se *SshConnectionError) Error() string {
	return se.errorMessage
}

func NewSshConnectionError(errorMessage string) *SshConnectionError {
	return &SshConnectionError{
		errorMessage: errorMessage,
	}
}

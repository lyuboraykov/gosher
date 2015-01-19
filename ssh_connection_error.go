package gosher

type SshConnectionError struct {
	errorMessage string
}

func (se *SshConnectionError) Error() string {
	return se.errorMessage
}

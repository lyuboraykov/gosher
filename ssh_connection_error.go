package gosher

type sshConnectionError struct {
	errorMessage string
}

func (se *sshConnectionError) Error() string {
	return se.errorMessage
}

package gosher

type sshResponse struct {
	exitCode int
	stdOut   bytes.Buffer
	stdErr   bytes.Buffer
}

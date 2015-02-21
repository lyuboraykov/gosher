package gosher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var isServerStarted bool

func initServer() {
	if !isServerStarted {
		started := make(chan bool)
		go runTestServer(started)
		isServerStarted = <-started
	}
}

func TestRun(t *testing.T) {
	initServer()
	s, clientError := NewSshClient("localhost", "foo", PasswordAuthentication, "bar")
	s.Port = 2200
	assert.Nil(t, clientError, "NewSshClient returned an error")
	assert.NotNil(t, s, "NewSshClient returned a nil client")
	response, runError := s.Run("echo Hello")
	fmt.Println("here")
	assert.Nil(t, runError, "Run returned an error")
	assert.NotNil(t, response, "Run returned a nil response")
	responseStdOut := response.StdOut.String()
	assert.Equal(t, "Hello", responseStdOut, "Different standard output, expected Hello, actual was: "+responseStdOut)
}

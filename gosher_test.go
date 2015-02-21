package gosher

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

const (
	serverAddress = "shell.xshellz.com"
	username      = "lyuboraykov"
	password      = "164o7w3sld"
)

func TestRun(t *testing.T) {
	s, clientError := NewSshClient(serverAddress, username, PasswordAuthentication, password)
	assert.Nil(t, clientError, "NewSshClient returned an error")
	assert.NotNil(t, s, "NewSshClient returned a nil client")
	response, runError := s.Run("echo Hello")
	assert.Nil(t, runError, "Run returned an error")
	assert.NotNil(t, response, "Run returned a nil response")
	responseStdOut := response.StdOut.String()
	assert.Equal(t, "Hello\n", responseStdOut, "Different standard output, expected Hello, actual was: "+responseStdOut)
}

func TestUpload(t *testing.T) {
	s, clientError := NewSshClient(serverAddress, username, PasswordAuthentication, password)
	assert.Nil(t, clientError, "NewSshClient returned an error")
	assert.NotNil(t, s, "NewSshClient returned a nil client")
	ioutil.WriteFile("testfile", []byte("testcontent"), 777)
	response, uploadError := s.Upload("testfile", "testfile")
	assert.Nil(t, uploadError, "Run returned an error")
	assert.NotNil(t, response, "Run returned a nil response")
	catResponse, catError := s.Run("cat testfile")
	assert.Nil(t, catError, "reading test file returned and error")
	assert.NotNil(t, catResponse, "Cat run returned a nil response")
	responseStdOut := catResponse.StdOut.String()
	assert.Equal(t, "testcontent", responseStdOut, "Different standard output, expected testcontent, actual was: "+responseStdOut)
}

func TestDownload(t *testing.T) {
	s, clientError := NewSshClient(serverAddress, username, PasswordAuthentication, password)
	assert.Nil(t, clientError, "NewSshClient returned an error")
	assert.NotNil(t, s, "NewSshClient returned a nil client")
	echoResponse, echoError := s.Run("echo downloaded > downfile")
	assert.Nil(t, echoError, "creating test file returned and error")
	assert.NotNil(t, echoResponse, "Echo run returned a nil response")
	currentDirectory, dirErr := os.Getwd()
	assert.Nil(t, dirErr, "Error getting current directory")
	response, downloadError := s.Download("downfile", currentDirectory)
	assert.Nil(t, downloadError, "Run returned an error")
	assert.NotNil(t, response, "Run returned a nil response")
	content, readErr := ioutil.ReadFile("downfile")
	assert.Nil(t, readErr, "Reading the file returned an error")
	stringContent := string(content)
	assert.Equal(t, "downloaded", stringContent, "Different standard output, expected downloaded, actual was: "+stringContent)
}

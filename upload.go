package gosher

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (s *SshClient) uploadFile(localPath string, remotePath string) (*SshResponse, error) {
	response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)
	go func() {
		inPipe, _ := s.session.StdinPipe()
		defer inPipe.Close()
		writeFileInPipe(inPipe, localPath, filepath.Base(remotePath))
	}()

	if err := s.session.Run("/usr/bin/scp -qvrt " + filepath.Dir(remotePath)); err != nil {
		return response, NewSshConnectionError("There was an error while uploading: " + err.Error())
	}
	return response, nil
}

func (s *SshClient) uploadFolder(localPath string, remotePath string) (*SshResponse, error) {
	response := NewSshResponse(s.address, s.session.Stdout, s.session.Stderr)
	go func() {
		inPipe, _ := s.session.StdinPipe()
		defer inPipe.Close()
		fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FOLDER, filepath.Base(remotePath))
		writeDirectoryContents(inPipe, localPath)
		fmt.Fprintln(inPipe, SCP_PUSH_END_FOLDER)
	}()

	if err := s.session.Run("/usr/bin/scp -qvrt " + filepath.Dir(remotePath)); err != nil {
		return response, NewSshConnectionError("Error while uploading: " + err.Error())
	}
	return response, nil
}

func writeDirectoryContents(inPipe io.WriteCloser, dir string) {
	fi, _ := ioutil.ReadDir(dir)
	for _, f := range fi {
		if f.IsDir() {
			fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FOLDER, f.Name())
			writeDirectoryContents(inPipe, dir+"/"+f.Name())
			fmt.Fprintln(inPipe, SCP_PUSH_END_FOLDER)
		} else {
			writeFileInPipe(inPipe, dir+"/"+f.Name(), f.Name())
		}
	}
}

func writeFileInPipe(inPipe io.WriteCloser, src string, remoteName string) {
	fileSrc, _ := os.Open(src)
	//Get file size
	srcStat, _ := fileSrc.Stat()
	// Print the file content
	fmt.Fprintln(inPipe, SCP_PUSH_BEGIN_FILE, srcStat.Size(), remoteName)
	io.Copy(inPipe, fileSrc)
	fmt.Fprint(inPipe, SCP_PUSH_END)
}

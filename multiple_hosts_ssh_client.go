package gosher

import "strconv"

type MultipleHostsSshClient struct {
	Hosts []*Host
}

// Constructor method for MultipleHostsSshClient
// Use this client when dealing with multiple hosts asynchronously.
func NewMultipleHostsSshClient(hosts ...*Host) *MultipleHostsSshClient {
	return &MultipleHostsSshClient{
		Hosts: hosts,
	}
}

// Executes shell command on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels
func (msc *MultipleHostsSshClient) Run(command string) {
	for i := range msc.Hosts {
		go func(i int) {
			result, err := msc.Hosts[i].Client.Run(command)
			if err != nil {
				msc.Hosts[i].ErrorChannel <- err
				return
			}
			msc.Hosts[i].ResultChannel <- result
		}(i)
	}
}

// Executes shell script on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels
func (msc *MultipleHostsSshClient) RunScript(filePath string) {
	for i := range msc.Hosts {
		go func(i int) {
			result, err := msc.Hosts[i].Client.RunScript(filePath)
			if err != nil {
				msc.Hosts[i].ErrorChannel <- err
				return
			}
			msc.Hosts[i].ResultChannel <- result
		}(i)
	}
}

// Uploads a file/folder to all hosts of the MultipleHostsSshClient.
// The sshResponse is passed via the channels of the hosts
func (msc *MultipleHostsSshClient) Upload(localPath string, remotePath string) {
	for i := range msc.Hosts {
		go func(i int) {
			result, err := msc.Hosts[i].Client.Upload(localPath, remotePath)
			if err != nil {
				msc.Hosts[i].ErrorChannel <- err
				return
			}
			msc.Hosts[i].ResultChannel <- result
		}(i)
	}
}

// Downloads files/folders from all hosts of the MultipleHostsSshClient's list.
// They will be suffixed with the index of the host they are downloaded from
func (msc *MultipleHostsSshClient) Download(remotePath string, localPath string) {
	for i, host := range msc.Hosts {
		go func(i int) {
			suffixedDownloadPath := localPath + strconv.Itoa(i)
			result, err := msc.Hosts[i].Client.Download(remotePath, suffixedDownloadPath)
			if err != nil {
				host.ErrorChannel <- err
				return
			}
			host.ResultChannel <- result
		}(i)
	}
}

// Executes an function on a remote text file on all hosts.
// Can be used as an alternative of executing sed or awk on the remote machine.
// alterContentsFunction - the function to be executed, the content of the file as string will be
// passed to it and it should return the modified content.
func (msc *MultipleHostsSshClient) RunOnFile(filePath string,
	alterContentsFunction func(fileContent string) string) {
	for i := range msc.Hosts {
		go func(i int) {
			result, err := msc.Hosts[i].Client.RunOnFile(filePath, alterContentsFunction)
			if err != nil {
				msc.Hosts[i].ErrorChannel <- err
				return
			}
			msc.Hosts[i].ResultChannel <- result
		}(i)
	}
}

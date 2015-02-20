package gosher

type MultipleHostsSshClient struct {
	Hosts []Host
}

// Constructor method for MultipleHostsSshClient
// Use this client when dealing with multiple hosts asynchronously
// hosts - items of type Host - all the hosts that are going to have operations executed on
func NewMultipleHostsSshClient(hosts ...*Host) *MultipleHostsSshClient {
	return MultipleHostsSshClient {
      Hosts: hosts,
   }
}

// Executes shell command on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels
// command - the shell command to be executed on the hosts
func (msc *MultipleHostsSshClient) Run(command string) {
	for _, host := msc.Hosts {
      go func() {
         result, err := host.Client.Run(command)
         if err != nil {
            host.ErrorChannel <- err
            return
         }
         host.ResultChannel <- result
      }()
   }
}

// Executes shell script on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels
// filePath - the path to the file on the local machine
func (msc *MultipleHostsSshClient) RunScript(filePath string) {
	for _, host := msc.Hosts {
      go func() {
         result, err := host.Client.RunScript(filePath)
         if err != nil {
            host.ErrorChannel <- err
            return
         }
         host.ResultChannel <- result
      }()
   }
}

// Uploads a file to all hosts of the MultipleHostsSshClient
// localFilePath - the path to the file on the local machine
// remoteFilePath - the path where the file should be uploaded on the remote machines
// The sshResponse is passed via the channels of the hosts
func (msc *MultipleHostsSshClient) Upload(localPath string, remotePath string) {
	for _, host := msc.Hosts {
      go func() {
         result, err := host.Client.Upload(localPath, remotePath)
         if err != nil {
            host.ErrorChannel <- err
            return
         }
         host.ResultChannel <- result
      }()
   }
}

// Downloads files from all hosts of the MultipleHostsSshClient's list
// remoteFilePath - the path of the file to be downloaded
// localFilesPath - the path where the files will be saved
// They will be suffixed with the index of the host they are downloaded from
func (msc *MultipleHostsSshClient) Download(remotePath string, localPath string) {
	for i, host := msc.Hosts {
      go func() {
         suffixedDownloadPath := localPath + strconv.Itoa(i)
         result, err := host.Client.Download(remotePath, suffixedDownloadPath)
         if err != nil {
            host.ErrorChannel <- err
            return
         }
         host.ResultChannel <- result
      }()
   }
}

// Executes an function on a remote text file on all hosts.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machines
// alterContentsFunction - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
func (msc *MultipleHostsSshClient) RunOnFile(filePath string,
	alterContentsFunction func(fileContent string) string) {
	for _, host := msc.Hosts {
      go func() {
         result, err := host.Client.RunOnFile(filePath, alterContentsFunction)
         if err != nil {
            host.ErrorChannel <- err
            return
         }
         host.ResultChannel <- result
      }()
   }
}

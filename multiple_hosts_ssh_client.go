package gosher

type multipleHostsSshClient struct {
   hosts Host[]
}

// Constructor method for multipleHostsSshClient
// Use this client when dealing with multiple hosts asynchronously
// hosts - slice of items of type Host - all the hosts that are going
// to have operations executed on
func NewMultipleHostsSshClient(hosts *Host[]) (multipleHostsSshClient*, error) {

}

// Executes shell command on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels 
// command - the shell command to be executed on the hosts
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteCommandOnAllHosts(command string) error {

}

// Executes shell command only on selected hosts from the client's list
// The result from execution is passed via the hosts' channels
// command - the shell command to be executed on the hosts
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteCommandOnSelectedHosts(command string, hostIndexes ...int) error {

}

// Uploads a file to all hosts of the MultipleHostsSshClient
// localFilePath - the path to the file on the local machine
// remoteFilePath - the path where the file should be uploaded on the remote machines
// The sshResponse is passed via the channels of the hosts
// Returns an error if any has occured
func (msc *multipleHostsSshClient) UploadFileToAllHosts(localFilePath string, remoteFilePath string) error {}

// Uploads a file only to selected hosts of the MultipleHostsSshClient's list
// localFilePath - the path to the file on the local machine
// remoteFilePath - the path where the file should be uploaded on the remote machines
// The sshResponse is passed via the channels of the hosts
// Returns an error if any has occured
func (msc *multipleHostsSshClient) UploadFileToSelectedHosts(localFilePath string, remoteFilePath string, hostIndexes ...int) error {}

func (msc *multipleHostsSshClient) DownloadFileFromAllHosts(remoteFilePath string, 
      localFilePath string, renamingFunction(fileName string, host *Host) string) error {}

func (msc *multipleHostsSshClient) DownloadFileFromSelectedHosts(remoteFilePath string, 
      localFilePath string, renamingFunction(fileName string, host *Host) string, hostIndexes ...int) error {}

func (msc *multipleHostsSshClient) executeOnFileOnAllHosts(filePath string, alterContentsFunction func(fileContent string) string) error {}

func (msc *multipleHostsSshClient) executeOnFileOnSelectedHosts(filePath string, 
   alterContentsFunction func(fileContent string, hostIndexes ...int) string) error {}
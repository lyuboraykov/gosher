package gosher

type multipleHostsSshClient struct {
   hosts Host[]
}

// Constructor method for multipleHostsSshClient
// Use this client when dealing with multiple hosts asynchronously
// hosts - slice of items of type Host - all the hosts that are going
// to have operations executed on
func NewMultipleHostsSshClient(hosts ...*Host) (*multipleHostsSshClient, error) {

}

// Executes shell command on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels 
// command - the shell command to be executed on the hosts
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteCommandOnAllHosts(command string) error {

}

// Executes shell script only on selected hosts from the client's list
// The result from execution is passed via the hosts' channels
// command - the shell command to be executed on the hosts
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteCommandOnSelectedHosts(command string, hostIndexes ...int) error {

}

// Executes shell script on all hosts in a separate goroutine for each.
// The result from execution is passed via the hosts' channels 
// filePath - the path to the file on the local machine
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteScriptOnAllHosts(filePath string) error {

}

// Executes shell script only on selected hosts from the client's list
// The result from execution is passed via the hosts' channels
// filePath - the path to the file on the local machine
// Returns an error if any has occured.
func (msc *multipleHostsSshClient) ExecuteScriptOnSelectedHosts(filePath string, hostIndexes ...int) error {

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
// hostIndexes - the indexes of the hosts from the client's list 
// The sshResponse is passed via the channels of the hosts
// Returns an error if any has occured
func (msc *multipleHostsSshClient) UploadFileToSelectedHosts(localFilePath string, 
      remoteFilePath string, hostIndexes ...int) error {}

// Downloads files from all hosts of the MultipleHostsSshClient's list
// remoteFilePath - the path of the file to be downloaded
// localFilesPath - the path where the files will be saved
// They will be suffixed with the index of the host they are downloaded from
// Returns an error if one occurs.
func (msc *multipleHostsSshClient) DownloadFileFromAllHosts(remoteFilePath string, 
      localFilesPath string) error {}

// Downloads files only from selected hosts of the MultipleHostsSshClient's list
// remoteFilePath - the path of the file to be downloaded
// localFilesPath - the path where the files will be saved
// hostIndexes - the indexes of the hosts from the client's list
// They will be suffixed with the index of the host they are downloaded from
// Returns an error if one occurs.
func (msc *multipleHostsSshClient) DownloadFileFromSelectedHosts(remoteFilePath string, 
      localFilePath string, renamingFunction(fileName string, host *Host) string, hostIndexes ...int) error {}

// Executes an function on a remote text file on all hosts.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machines
// alterContentsFunction - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
// Returns an error if one occurs
func (msc *multipleHostsSshClient) ExecuteOnFileOnAllHosts(filePath string, 
      alterContentsFunction func(fileContent string) string) error {}

// Executes an function on a remote text file on all hosts.
// Can be used as an alternative of executing sed or awk on the remote machine.
// filePath - the path of the file on the remote machines
// alterContentsFunction - the function to be executed, the contents of the file as string will be
// passed to it and it should return the modified contents.
// hostIndexes - the indexes of the hosts from the client's list 
// Returns an error if one occurs.
func (msc *multipleHostsSshClient) ExecuteOnFileOnSelectedHosts(filePath string, 
      alterContentsFunction func(fileContent string) string, hostIndexes ...int) error {}
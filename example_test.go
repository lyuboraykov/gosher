package gosher_test

// import "fmt"

// import "gosher"

// func ExampleSynchronousHelloWorld() {
// 	client := NewSshClient("10.23.123.192", "root", PASSWORD_AUTH, "password")
// 	response, err := client.ExecuteCommand("echo 'Hello World!'")
// 	if err == nil {
// 		if response.ExitCode == 0 {
// 			fmt.Println(response.StdOut.String())
// 		}
//       else {
//          fmt.Printf("There was an error while executing the command: %s \n", response.StdErr.String())
//       }
// 	}
//    else {
//       fmt.Printf("There was an error while connecting to the server: %s \n", err.Error())
//    }
//    // Output: Hello World
// }

// func ExampleAsynchronousHelloWorld() {
// 	receiveChannel := make(chan<- *SshResponse)
// 	host1 = NewHost("10.23.123.191", "root", PASSWORD_AUTH, "password", receiveChannel)
// 	host2 = NewHost("10.23.123.192", "root", PASSWORD_AUTH, "password", receiveChannel)
// 	client = NewMultipleHostsSshClient(&host1, &host2)
// 	err := client.ExecuteCommandOnAllHosts("echo 'Hello World'")
// 	if err == nil {
// 		response1 := <-receiveChannel
// 		response2 := <-receiveChannel
// 		fmt.Println(response1.StdOut.String())
// 	}
//    // Output: Hello World
// }

// func ExampleSynchronousFileUploadWithKeyAuth() {
// 	client = NewSshClient("10.23.123.192", "root", KEY_AUTH, "~/.ssh/id_rsa.pub")
// 	response, err := client.UploadFile("./test_file", "/tmp/test_file")
//    if err == nil {
//       if response.ExitCode == 0 {
//          fmt.Println("File was uploaded successfully")
//       }
//    }
//    // Output: File was uploaded successfully
// }

// func ExampleExecuteOnFileOnSelectedHosts() {
//    receiveChannel := make(chan<- *SshResponse, 2)
//    host1 := NewHost("10.23.123.191", "root", PASSWORD_AUTH, "password", receiveChannel)
//    host2 := NewHost("10.23.123.192", "root", KEY_AUTH, "~/.ssh/id_rsa.pub", receiveChannel)
//    host3 := NewHost("10.23.123.193", "root", PASSWORD_AUTH, "password", receiveChannel)
//    client := NewMultipleHostsSshClient(host1, host2, host3)
//    err := client.ExecuteOnFileOnSelectedHosts("/tmp/test_file" func(fileContent string) {
//       return fileContent + " appended text"
//       }, 0, 1)
//    if err == nil {
//       response1 := <-receiveChannel
//       response2 := <-receiveChannel
//       fmt.Println(response1.ExitCode)
//    }
//    // Output: 0
// }

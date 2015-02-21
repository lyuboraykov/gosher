package gosher_test

import "fmt"

import "gosher"

func ExampleSynchronousHelloWorld() {
	client := NewSshClient("10.23.123.192", "root", gosher.PasswordAuthentication, "password")
	response, err := client.Run("echo 'Hello World!'")
	if err != nil {
		fmt.Printf("There was an error while executing the command: %s \n", response.StdErr.String())
	}
	fmt.Println(response.StdOut.String())
	// Output: Hello World
}

func ExampleAsynchronousHelloWorld() {
	receiveChannel := make(chan *SshResponse)
	errorChannel := mane(chan error)
	host1 = NewHost("10.23.123.191", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
	host2 = NewHost("10.23.123.192", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
	client = NewMultipleHostsSshClient(&host1, &host2)
	err := client.Run("echo 'Hello World'")
	if err == nil {
		response1 := <-receiveChannel
		response2 := <-receiveChannel
		fmt.Println(response1.StdOut.String())
	}
	// Output: Hello World
}

func ExampleSynchronousFileUploadWithKeyAuth() {
	client = NewSshClient("10.23.123.192", "root", KeyAuthentication, "~/.ssh/id_rsa.pub")
	response, err := client.Upload("./test_file", "/tmp/test_file")
	if err == nil {
		fmt.Println("File was uploaded successfully")
	}
	// Output: File was uploaded successfully
}

func ExampleRunOnFileOnMultipleHosts() {
	receiveChannel := make(chan *SshResponse, 2)
	errorChannel := make(chan error, 2)
	host1 := NewHost("10.23.123.191", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
	host2 := NewHost("10.23.123.192", "root", KeyAuthentication, "~/.ssh/id_rsa.pub", receiveChannel, errorChannel)
	host3 := NewHost("10.23.123.193", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
	client := NewMultipleHostsSshClient(host1, host2, host3)
	err := client.RunOnFile("/tmp/test_file", func(fileContent string) {
		return fileContent + " appended text"
	})
	if err == nil {
		response1 := <-receiveChannel
		response2 := <-receiveChannel
		fmt.Println("Success")
	}
	// Output: Success
}

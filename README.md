Gosher
======

Gosher is an SSH library for Go. 
It supports synchronous execution, file download/upload on single host 
and asynchronous operations on multiple hosts.

Installation
------------

To get gosher run: 
```shell
go get github.com/lyuboraykov/gosher
```
This will download, compile and install the package in your `$GOPATH` directory.
Alternatively, you can import the package like that:
```go
import "github.com/lyuboraykov/gosher"
```
and use `go get` without parameters.

Documentation
-----
You can read more about the interface here: [Godoc](https://godoc.org/github.com/lyuboraykov/gosher)

Usage
-----

There are two clients in the package one is SshClient, used for synchronous 
operations on a single host and the other is MultipleHostsSshClient for 
async operations on multiple hosts.

Here is an example synchronous Hello World on a single host: 
```go
import "github.com/lyuboraykov/gosher"

// ...
client := gosher.NewSshClient("10.23.123.192", "root", gosher.PasswordAuthentication, "password")
response, err := client.Run("echo 'Hello World!'")
if err == nil {
   fmt.Println(response.StdOut.String())
}
else {
   fmt.Printf("There was an error while connecting to the server: %s \n", err.Error())
}
```

Here is a simple file upload: 
```go
import "github.com/lyuboraykov/gosher"

// ...
client := gosher.NewSshClient("10.23.123.192", "root", gosher.KeyAuthentication, "/home/user/.ssh/id_rsa")
response, err := client.Upload("somefile", "somelocation")
if err == nil {
   fmt.Println(response.StdOut.String())
}
else {
   fmt.Printf("There was an error while connecting to the server: %s \n", err.Error())
}
```

And here is a hello world on two hosts async: 

```go
receiveChannel := make(chan *SshResponse)
errorChannel := make(chan error)
host1 = gosher.NewHost("10.23.123.191", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
host2 = gosher.NewHost("10.23.123.192", "root", PasswordAuthentication, "password", receiveChannel, errorChannel)
client = NewMultipleHostsSshClient(&host1, &host2)
err := client.Run("echo 'Hello World'")
if err == nil {
   response1 := <-receiveChannel
   response2 := <-receiveChannel
   fmt.Println(response1.StdOut.String())
}
```

Now let's get more advanced and execute a function on a file on multiple hosts:

```go
receiveChannel := make(chan *SshResponse, 2)
errorChannel := make(chan error, 2)
host1 := NewHost("10.23.123.191", "root", gosher.PasswordAuthentication, "password", receiveChannel, errorChannel)
host2 := NewHost("10.23.123.192", "root", goser.KeyAuthentication, "~/.ssh/id_rsa", receiveChannel, errorChannel)
host3 := NewHost("10.23.123.193", "root", gosher.PasswordAuthentication, "password", receiveChannel, errorChannel)
client := NewMultipleHostsSshClient(host1, host2, host3)
err := client.RunOnFile("/tmp/test_file" func(fileContent string) {
   return fileContent + " appended text"
   })
if err == nil {
   response1 := <-receiveChannel
   response2 := <-receiveChannel
   fmt.Println("Success!")
}
```

Features
--------
All of these features are supported both on a single host and on multiple hosts

*   **Execute Command** Executes a simple shell command

*   **Execute Script** Executes a local script on remote machine

*   **Upload/Download** provides scp functionality

*   **Execute on file** executes a function on a remote file, can be used
    instead of awk/sed

Todo
----

Do a full test coverage.

License
-------
This package is distributed under the MIT License:

```
The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
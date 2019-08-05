package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

var socksVersion byte = 05

func main() {
	fmt.Println(socksVersion)
	// https://rushter.com/blog/python-socks-server/
	server, err := net.Listen("tcp", ":3030")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := server.Accept()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := make([]byte, 257)
	reqLen, err := conn.Read(buf)
	fmt.Println(buf)
	checkHeaders(buf)
	fmt.Println(reqLen)
	// Clients sends the following msg in that format ... ->
	// version	nmethods	methods
	// 1 byte	1 byte	0 to 255 bytes
	// version is in our case is 5.
	// nmethods - number of auth methods.
	// Let's implement the basic usr/pass method. Which is according to the protocol is 02. (2)
	// If 02 method is in possible methods, we choose that and send it back to client.
	// Then, we wait for the followings message format:
	// version	ulen	uname	plen	passwd
	// 1 byte	1 byte	0 to 255 bytes	1 byte	0 to 255 bytes
	// ulen - length of the username, plen - length of password. Helps with parsing.
	// If usr/pass checks out - we send the status 0 to the client which indicates successs.
	// Now we ready to accept requests in the following format:
	// version	cmd	rsv	atyp	dst.addr	dst.port
	// 1 byte	1 byte	1 byte	1 byte	4 to 255 bytes	2 bytes
	// Process the requests and send back the response to the client in the following format:
	// version	rep	rsv	atyp	bnd.addr	bnd.port
	// 1 byte	1 byte	1 byte	1 byte	4 to 255 bytes	2 bytes
}

func checkHeaders(buf []byte) {
	bufferReader := bytes.NewReader(buf)

	version, err := bufferReader.ReadByte()
	if err != nil {
		fmt.Println("Problems reading headers.")
		os.Exit(1)
	}
	if version == socksVersion {
		fmt.Println("Socks5 version verified.")
	}
	nmethods, err := bufferReader.ReadByte()
	if err != nil {
		fmt.Println("Problems reading headers.")
		os.Exit(1)
	}
	var methods = []byte{}
	for i := 1; i <= int(nmethods); i++ {

	}
}

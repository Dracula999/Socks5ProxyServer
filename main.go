package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

var socksVersion byte = 05
var USERNAME = "usr"
var PASSWORD = "pass"

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
	conn.Read(buf)
	// checking headers and auth type
	auth := checkHeaders(buf)
	if auth == true {
		conn.Write([]byte{5, 2})
	}
	next := make([]byte, 257)
	conn.Read(next)
	// parse auth creds
	username, password := parseAuthCredentials(next)
	fmt.Println(username, password)
	// check auth creds
	if (USERNAME == username) && (PASSWORD == password) {
		fmt.Println("Client has given right credentials.")
		conn.Write([]byte{1, 0})
	} else {
		fmt.Println("Client has given wrong credentials.")
		conn.Write([]byte{5, 1})
	}
	reqBytes := make([]byte, 256)
	conn.Read(reqBytes)
	fmt.Println(reqBytes)
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
func parseAuthCredentials(buf []byte) (string, string) {
	bufferReader := bytes.NewReader(buf)
	bufferReader.ReadByte()
	usernameLen, _ := bufferReader.ReadByte()
	username := make([]byte, int(usernameLen))
	bufferReader.Read(username)
	passwordLen, _ := bufferReader.ReadByte()
	password := make([]byte, int(passwordLen))
	bufferReader.Read(password)
	return string(username), string(password)
}
func checkHeaders(buf []byte) bool {
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
	n := int(nmethods)
	var methods = make([]byte, n)
	bufferReader.Read(methods)
	fmt.Print("Methods: ")
	fmt.Println(methods)
	for _, method := range methods {
		if method == 2 {
			return true
		}
	}
	return false
}

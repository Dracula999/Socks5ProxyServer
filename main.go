package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
)

var socksVersion byte = 5
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
	addr := parseAddr(reqBytes)
	dialAndGetAddr(addr)
	// Process the requests and send back the response to the client in the following format:
	// version	rep	rsv	atyp	bnd.addr	bnd.port
	// 1 byte	1 byte	1 byte	1 byte	4 to 255 bytes	2 bytes

}

func dialAndGetAddr(addr string) {
	remoteConn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Couldn't connect to that addr " + addr)
	} else {
		fmt.Println("Connected.")
	}
	ip, port, err := net.SplitHostPort(remoteConn.LocalAddr().String())
	fmt.Println(port)
	destAddr := net.ParseIP(ip).To4()
	fmt.Println(destAddr)
	val, _ := strconv.ParseUint(port, 16, 2)
	fmt.Println(val)
}

func parseAddr(buf []byte) string {
	// Now we ready to accept requests in the following format:
	// version	cmd	rsv	atyp	dst.addr	dst.port
	// 1 byte	1 byte	1 byte	1 byte	4 to 255 bytes	2 bytes
	bufferReader := bytes.NewReader(buf)
	bufferReader.ReadByte()
	connectionType, _ := bufferReader.ReadByte()
	if connectionType != 1 {
		return ""
	}
	// for now we implement only first type which is to connect.
	// rsv we just skip that byte. Idk why it's reserved.
	bufferReader.ReadByte()
	// atype is for the addr type which can be ipv4, domain name or ipv6. For now let's implement ipv4
	atyp, _ := bufferReader.ReadByte()
	if atyp != 1 {
		return ""
	}
	// now we read addr and port.
	addr := make([]byte, 4)
	port := make([]byte, 2)

	bufferReader.Read(addr)
	bufferReader.Read(port)
	fmt.Println(port)
	ip := net.IP(addr)
	portInt := binary.BigEndian.Uint16(port)
	fullAddr := ip.String() + ":" + string(strconv.Itoa(int(portInt)))
	fmt.Println(fullAddr)
	return fullAddr
}
func parseAuthCredentials(buf []byte) (string, string) {
	// Then, we wait for the followings message format:
	// version	ulen	uname	plen	passwd
	// 1 byte	1 byte	0 to 255 bytes	1 byte	0 to 255 bytes
	// ulen - length of the username, plen - length of password. Helps with parsing.
	// If usr/pass checks out - we send the status 0 to the client which indicates successs.

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
	// Clients sends the following msg in that format ... ->
	// version	nmethods	methods
	// 1 byte	1 byte	0 to 255 bytes
	// version is in our case is 5.
	// nmethods - number of auth methods.
	// Let's implement the basic usr/pass method. Which is according to the protocol is 02. (2)
	// If 02 method is in possible methods, we choose that and send it back to client.

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

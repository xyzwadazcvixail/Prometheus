package main

import (
	"bufio"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var encode_base64 = func(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

var decode_base64 = func(data []byte) string {
	decoded, err := b64.StdEncoding.DecodeString(string(data))
	check(err)
	return string(decoded)
}

func process_flags() (string, int, string) {
	var (
		port    int
		address string
		protocol string
	)
	flag.IntVar(&port, "port", 8080, "Port Number")
	flag.StringVar(&address, "address", "127.0.0.1", "IPv4")
	flag.StringVar(&protocol, "protocol", "udp", "Protocol Name")
	flag.Parse()
	return address, port, protocol
}

func main() {
	buf := make([]byte, 262_144)
	address, port, protocol := process_flags()
	remoteAddress := fmt.Sprintf("%s:%d", address, port)
	conn, err := net.Dial(protocol, remoteAddress)
	check(err)
	defer conn.Close()
	for {
		fmt.Printf("Enter command {%s//%s:%d}:\n", protocol, address, port)
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		check(err)
		encoded := encode_base64(command)
		fmt.Fprintf(conn, encoded)
		n, err := bufio.NewReader(conn).Read(buf)
		fmt.Println(decode_base64(buf[:n]))
		check(err)
	}
}

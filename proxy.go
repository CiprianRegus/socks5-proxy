package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const VER_5 = 0x05
const NO_AUTH = 0x00
const GSSAPI = 0x01
const USER_PASS = 0x02
const NO_METHOD = 0xff

func err(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func get_addr_port(request []byte) ([]byte, []byte) {

	/*
	   +----+-----+-------+------+----------+----------+
	   |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	   +----+-----+-------+------+----------+----------+
	   | 1  |  1  | X'00' |  1   | Variable |    2     |
	   +----+-----+-------+------+----------+----------+
	*/

	addr := make([]byte, len(request)-6)
	port := make([]byte, 2)

	addr = request[5 : len(request)-2]
	port = request[len(request)-2:]

	return addr, port
}

func handle_new_connection(conn net.Conn) {

	//defer conn.Close()

	msg := make([]byte, 100)
	_, e := conn.Read(msg)
	err(e)
	//fmt.Println(msg)

	if msg[0] != VER_5 {
		panic("Invalid version")
	}
	method := msg[2]

	switch method {
	case NO_AUTH:
		{
			estable_connection(conn, NO_AUTH)
		}
	default:
		panic("Method not supported")
	}
}

func estable_connection(conn net.Conn, method int) {

	conn.Write([]byte{VER_5, byte(method)})
	request := make([]byte, 100)
	req_len, _ := conn.Read(request)

	ln, e := net.Listen("tcp", "localhost:42412")
	err(e)

	for {
		conn, e := ln.Accept()
		err(e)
		go func(conn net.Conn) {
			fmt.Println("New connection")
		}(conn)
	}

	addr, port := get_addr_port(request[:req_len])
	p := []byte{0xac, 0xa5}
	conn.Write(append(append([]byte{VER_5, 0x0, 0x0, 0x3}, []byte("localhost")...), p...))

	/*
		dst_conn, _ := net.Dial("tcp", "www.google.com:8080")
		fmt.Fprintf(dst_conn, "GET / HTTP/1.0\r\n\r\n")
	*/

	fmt.Println(addr, string(addr), binary.BigEndian.Uint16(port))
	remote_addr := string(addr) + ":" + "80"
	res := make([]byte, 1000000)
	conn2, e := net.DialTimeout("tcp", remote_addr, 3*time.Second)
	//fmt.Println(request)
	err(e)
	fmt.Fprintf(conn2, "GET / HTTP/1.0\r\n\r\n")

	_, readerr := conn2.Read(res)
	err(readerr)
	conn.Write(res)
	fmt.Println(string(res))
}

func main() {

	ln, e := net.Listen("tcp", "localhost:1080")
	err(e)

	for {
		conn, e := ln.Accept()
		err(e)
		go handle_new_connection(conn)
	}
}

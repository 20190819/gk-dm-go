package net

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestNetServer(t *testing.T) {

	Listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}

	for {
		conn, err := Listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			Handle(conn)
		}()
	}
}

func Handle(conn net.Conn) {

	for {
		msgHead := make([]byte, 8)
		_, err := conn.Read(msgHead)
		if err != nil {
			conn.Close()
			return
		}
		msgLen:=binary.BigEndian.Uint64(msgHead)
		resp:=make([]byte,msgLen)
		_, err = conn.Read(resp)
		if err != nil {
			conn.Close()
			return
		}

		_, err = conn.Write([]byte("hello world"))
		if err != nil {
			conn.Close()
			return
		}
	}

}

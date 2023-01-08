package net

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestNetClient(t *testing.T) {

	conn, err := net.Dial("tcp", ":8082")
	if err != nil {
		t.Fatal(err)
	}

	msg := "how are you golang"
	msglen := len(msg)
	msgHead := make([]byte, 8)
	binary.BigEndian.PutUint64(msgHead, uint64(msglen))
	data := append(msgHead, []byte(msg)...)
	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		return
	}

	respbs := make([]byte, 16)
	_, err = conn.Read(respbs)
	if err != nil {
		conn.Close()
		return
	}

}

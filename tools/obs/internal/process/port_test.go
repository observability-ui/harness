package process_test

import (
	"net"
	"testing"

	"obs/internal/process"
)

func TestCheckPort_Available(t *testing.T) {
	if err := process.CheckPort(59123); err != nil {
		t.Fatalf("unused port should be available: %v", err)
	}
}

func TestCheckPort_InUse(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	if err := process.CheckPort(port); err == nil {
		t.Fatal("occupied port should return error")
	}
}

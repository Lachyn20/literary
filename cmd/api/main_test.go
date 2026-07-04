package main

import (
	"net"
	"strconv"
	"testing"
)

func TestListenWithFallbackUsesAlternatePortWhenPrimaryIsBusy(t *testing.T) {
	busy, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to reserve a test port: %v", err)
	}
	defer busy.Close()

	port := strconv.Itoa(busy.Addr().(*net.TCPAddr).Port)
	listener, actualAddr, err := listenWithFallback(port)
	if err != nil {
		t.Fatalf("expected fallback listener, got error: %v", err)
	}
	defer listener.Close()

	if listener.Addr().(*net.TCPAddr).Port <= 0 {
		t.Fatalf("expected a valid listening port")
	}

	if actualAddr == "" {
		t.Fatalf("expected non-empty listening address")
	}
}

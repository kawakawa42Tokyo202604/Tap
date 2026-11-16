package main

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func receiveMessage(t *testing.T, client *protocolClient) serverMessage {
	t.Helper()
	select {
	case message, open := <-client.incoming:
		if !open {
			t.Fatal("connection closed before the expected message")
		}
		return message
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for a protocol message")
		return serverMessage{}
	}
}

func TestProtocolGreetingResponseAndEvent(t *testing.T) {
	local, remote := net.Pipe()
	client := newProtocolClient(local)
	defer client.Close()
	go func() {
		defer remote.Close()
		fmt.Fprintln(remote, "OK hello proto=1")
		scanner := bufio.NewScanner(remote)
		if scanner.Scan() {
			fmt.Fprintln(remote, `OK {"hp":10}`)
			fmt.Fprintln(remote, "EVT room_changed")
		}
	}()
	if got := receiveMessage(t, client); got.kind != messageGreeting {
		t.Fatalf("greeting kind: %v", got.kind)
	}
	if err := client.Send("STATUS"); err != nil {
		t.Fatal(err)
	}
	if got := receiveMessage(t, client); got.kind != messageResponse || got.command != "STATUS" || got.request != "STATUS" {
		t.Fatalf("response not correlated to command: %+v", got)
	}
	if got := receiveMessage(t, client); got.kind != messageEvent {
		t.Fatalf("event kind: %v", got.kind)
	}
}

func TestProtocolRejectsInvalidFraming(t *testing.T) {
	local, remote := net.Pipe()
	defer remote.Close()
	client := newProtocolClient(local)
	defer client.Close()
	for _, command := range []string{"", " ", "LOOK\nQUIT", "LOOK\rQUIT", string([]byte{0xff})} {
		if client.Send(command) == nil {
			t.Errorf("accepted invalid command %q", command)
		}
	}
	client.Close()
	if client.Send("LOOK") == nil {
		t.Fatal("accepted a command after Close")
	}
}

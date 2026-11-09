package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

type testConnection struct {
	net.Conn
}

func (connection testConnection) CloseWrite() error {
	return nil
}

func TestRawRelayAndQuit(t *testing.T) {
	for _, retryQuit := range []bool{false, true} {
		t.Run(fmt.Sprintf("retry=%t", retryQuit), func(t *testing.T) {
			server, client := net.Pipe()
			defer server.Close()
			defer client.Close()
			deadline := time.Now().Add(5 * time.Second)
			if err := server.SetDeadline(deadline); err != nil {
				t.Fatal(err)
			}
			if err := client.SetDeadline(deadline); err != nil {
				t.Fatal(err)
			}
			commands := []string{"CONNECT アリス", "QUIT"}
			responses := []string{"OK connected", "OK bye"}
			if retryQuit {
				commands = append(commands, "WHO", "QUIT")
				responses = []string{"OK connected", "ERR 500 STATE_ERROR", "OK players=1", "OK bye"}
			}
			result := make(chan error, 1)
			go func() {
				if _, err := fmt.Fprintln(server, "OK hello proto=1"); err != nil {
					result <- err
					return
				}
				reader := bufio.NewScanner(server)
				for index, want := range commands {
					if !reader.Scan() {
						result <- fmt.Errorf("missing command %q: %v", want, reader.Err())
						return
					}
					if got := reader.Text(); got != want {
						result <- fmt.Errorf("command = %q, want %q", got, want)
						return
					}
					if _, err := fmt.Fprintf(server, "EVT STATS players=1\n%s\n", responses[index]); err != nil {
						result <- err
						return
					}
				}
				result <- nil
			}()
			runClient(testConnection{client}, strings.NewReader("\n"+strings.Join(commands, "\n")+"\n"))
			if err := <-result; err != nil {
				t.Fatal(err)
			}
		})
	}
}

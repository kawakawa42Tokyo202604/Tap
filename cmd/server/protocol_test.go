package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestInvalidCommandsAndArguments(t *testing.T) {
	server := testServer(t)
	client := startTestClient(t, server)
	for _, command := range []string{
		"FROB", "CONNECT", "CONNECT a b", "CONNECT bad\x00name",
		"CONNECT \xff", "LANG fr", "LOOK", "MOVE east", "TAKE item.none",
		"DROP item.none", "INVENTORY", "STATUS", "TALK npc.none",
		"ATTACK npc.none", "QUEST npc.none", "QUESTS", "CHAT GLOBAL hello",
		"GROUP CREATE", "WHO extra", "QUIT extra",
	} {
		t.Run(fmt.Sprintf("before_%q", command), func(t *testing.T) {
			client.command(t, command, "ERR 400 BAD_REQUEST")
		})
	}
	client.command(t, "LANG ja", "OK lang=ja")
	client.command(t, "CONNECT アリス", "OK connected")
	for _, command := range []string{
		"CONNECT bob", "LANG en", "MOVE", "MOVE east extra", "LOOK extra",
		"CHAT UNKNOWN hello", "CHAT GLOBAL", "CHAT GLOBAL bad\x00message",
		"GROUP UNKNOWN", "GROUP INVITE", "GROUP CREATE extra",
	} {
		t.Run(fmt.Sprintf("after_%q", command), func(t *testing.T) {
			client.command(t, command, "ERR 400 BAD_REQUEST")
		})
	}
	client.command(t, "MOVE nowhere", "ERR 301 NO_EXIT")
	client.command(t, "CHAT GROUP hello", "ERR 401 NOT_IN_GROUP")
	client.command(t, "WHO", "OK players=1")
	client.command(t, "QUIT", "OK bye")
}

func TestSplitAndCoalescedCommands(t *testing.T) {
	client := startTestClient(t, testServer(t))
	writeDone := make(chan error, 1)
	go func() {
		for _, chunk := range []string{"CON", "NECT alice\nWHO\n", "\n \nWHO\r\n"} {
			if _, err := fmt.Fprint(client.conn, chunk); err != nil {
				writeDone <- err
				return
			}
		}
		writeDone <- nil
	}()
	client.expect(t, "OK connected")
	client.expect(t, "OK players=1")
	client.expect(t, "OK players=1")
	if err := <-writeDone; err != nil {
		t.Fatal(err)
	}
	client.command(t, "QUIT", "OK bye")
}

func TestLookAndChat(t *testing.T) {
	server := testServer(t)
	client := startTestClient(t, server)
	client.command(t, "CONNECT alice", "OK connected")
	if _, err := fmt.Fprintln(client.conn, "LOOK"); err != nil {
		t.Fatal(err)
	}
	response := client.response(t)
	if !strings.HasPrefix(response, "OK ") {
		t.Fatalf("LOOK response = %q", response)
	}
	var view struct {
		Room    roomView `json:"room"`
		Players []string `json:"players"`
	}
	if err := json.Unmarshal([]byte(strings.TrimPrefix(response, "OK ")), &view); err != nil {
		t.Fatal(err)
	}
	if view.Room.ID != server.world.StartRoomID || len(view.Players) != 1 || view.Players[0] != "alice" {
		t.Fatalf("unexpected LOOK contents: %+v", view)
	}
	client.command(t, "CHAT GLOBAL こんにちは  world", "OK")
	line, err := client.reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if want := "EVT GLOBAL CHAT alice こんにちは  world\n"; line != want {
		t.Fatalf("chat event = %q, want %q", line, want)
	}
	client.command(t, "QUIT", "OK bye")
}

func TestAbruptDisconnectPersistsPlayer(t *testing.T) {
	server := testServer(t)
	client := startTestClient(t, server)
	client.command(t, "CONNECT alice", "OK connected")
	client.conn.Close()
	select {
	case <-client.done:
	case <-time.After(5 * time.Second):
		t.Fatal("server did not finish disconnected client")
	}
	server.mu.Lock()
	remaining := len(server.players)
	server.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("online players after disconnect = %d", remaining)
	}
	players, err := server.loadPlayers()
	if err != nil || players["alice"] == nil {
		t.Fatalf("saved players = %v, error = %v", players, err)
	}
}

func TestInvalidWorldFiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "world.json")
	if _, err := loadWorld(path); err == nil {
		t.Fatal("missing world file was accepted")
	}
	for _, data := range []string{
		"{", "null", "{}",
		`{"start_room_id":"x","rooms":{"x":{"id":"x","exits":{"east":"missing"}}}}`,
	} {
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadWorld(path); err == nil {
			t.Errorf("invalid world was accepted: %s", data)
		}
	}
}

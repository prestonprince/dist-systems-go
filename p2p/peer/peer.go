package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var DIRECTORY_URL = "http://localhost:8080"

type Peer struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type Message struct {
	From    string `json:"from"`
	Content string `json:"content"`
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run peer.go <peer_id> <peer_addr>")
	}

	peerID := os.Args[1]
	peerAddr := os.Args[2]

	peer := Peer{ID: peerID, Addr: peerAddr}
	registerPeer(peer)

	http.HandleFunc("/message", handleMessage)
	go func() {
		log.Fatal(http.ListenAndServe(peerAddr, nil))
	}()

	for {
		peers := getPeer()
		for _, p := range peers {
			if p.ID != peerID {
				sendMessage(p, Message{
					From:    peerID,
					Content: fmt.Sprintf("hello from %s\n", peer.ID),
				})
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func registerPeer(peer Peer) {
	body, _ := json.Marshal(peer)
	resp, err := http.Post(
		DIRECTORY_URL+"/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Printf("failed to register: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func getPeer() []Peer {
	resp, err := http.Get(DIRECTORY_URL + "/peers")
	if err != nil {
		log.Printf("failed to get peers: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var peers []Peer
	json.NewDecoder(resp.Body).Decode(&peers)
	return peers
}

func sendMessage(peer Peer, msg Message) {
	body, _ := json.Marshal(msg)
	resp, err := http.Post(
		"http://"+peer.Addr+"/message",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Printf("failed to send message: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("received message from %s: %s\n", msg.From, msg.Content)
}

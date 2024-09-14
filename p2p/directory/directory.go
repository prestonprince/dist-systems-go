package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Peer struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type Directory struct {
	peers map[string]Peer
	mu    sync.RWMutex
}

func NewDirectory() *Directory {
	return &Directory{
		peers: make(map[string]Peer),
	}
}

func (d *Directory) RegisterPeer(w http.ResponseWriter, r *http.Request) {
	var peer Peer
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	d.mu.Lock()
	d.peers[peer.ID] = peer
	d.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (d *Directory) GetPeers(w http.ResponseWriter, r *http.Request) {
	d.mu.RLock()
	peers := make([]Peer, 0, len(d.peers))
	for _, peer := range d.peers {
		peers = append(peers, peer)
	}
	d.mu.RUnlock()

	json.NewEncoder(w).Encode(peers)
}

func main() {
	dir := NewDirectory()

	http.HandleFunc("/register", dir.RegisterPeer)
	http.HandleFunc("/peers", dir.GetPeers)

	fmt.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

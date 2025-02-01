package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type File struct {
	Name        string
	Size        int
	Uploaded_by string
	Date        int
}

type FileRegistry struct {
	files map[File]*net.TCPConn
	mu    sync.Mutex
}

const (
	broadcastPort = 42069
	broadcastIP   = "127.255.255.255" // change to 255 in prod
)

func NewFileRegistry() *FileRegistry {
	return &FileRegistry{
		files: make(map[File]*net.TCPConn, 0),
	}
}

// Insert performs a thread safe append to files slice.
func (fr *FileRegistry) Insert(item File, conn *net.UDPConn, tcpConn *net.TCPConn) {
	fr.mu.Lock()
	fr.files[item] = tcpConn
	fr.mu.Unlock()
	fr.SyncWrite(item, conn)
}

// SyncWrite broadcasts the newly uploaded File item.
//
// Should the caller Close the connection after successful broadcast ???
func (fr *FileRegistry) SyncWrite(item File, conn *net.UDPConn) {
	broadcastAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastIP, broadcastPort))
	jsonF, _ := json.Marshal(item)
	_, err := conn.WriteToUDP(jsonF, broadcastAddr)
	if err != nil {
		log.Println("Sync Write: Error broadcasting file for sync: ", err)
	}
}

// SyncRead listens for File update broadcasts and updates the local File copy.
func (fr *FileRegistry) SyncRead(tcpConn *net.TCPConn) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: broadcastPort})
	if err != nil {
		log.Println("Sync Read: Error connecting to broadcast: ", err)
		return
	}

	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("Sync Read: Error getting file part: ", err)
	}

	jsonF := buff[:n]
	var file File

	err = json.Unmarshal(jsonF, &file)
	if err != nil {
		log.Println("Sync Read: Error unmarshaling to json: ", err)
	}

	fr.mu.Lock()
	fr.files[file] = tcpConn
	fr.mu.Unlock()
}

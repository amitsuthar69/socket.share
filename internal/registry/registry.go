package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"socket-share/internal/discovery"
	"sync"
)

type File struct {
	Name        string
	Size        int
	Uploaded_by string
	Date        int
}

type FileRegistry struct {
	files map[File]string
	mu    sync.Mutex
}

const (
	broadcastPort = 42069
	broadcastIP   = "127.255.255.255" // change to 255 in prod
)

func NewFileRegistry() *FileRegistry {
	return &FileRegistry{
		files: make(map[File]string, 0),
	}
}

// Insert performs a thread safe append to files map.
func (fr *FileRegistry) Insert(item File, conn *net.UDPConn, ip string) {
	fr.mu.Lock()
	fr.files[item] = ip
	fr.mu.Unlock()
	fr.SyncWrite(item, conn)
}

// SyncWrite broadcasts the newly uploaded File item.
//
// Should the caller Close the connection after successful broadcast ???
func (fr *FileRegistry) SyncWrite(item File, conn *net.UDPConn) {
	broadcastAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastIP, broadcastPort))

	jsonF, err := json.Marshal(item)
	if err != nil {
		log.Println("Sync Write: Error marshaling to json: ", err)
	}

	n, err := conn.WriteToUDP(jsonF, broadcastAddr)
	if err != nil {
		log.Println("Sync Write: Error broadcasting file for sync: ", err)
	}

	log.Printf("Sync Write: Written %d bytes to sync.", n)
}

// SyncRead listens for File update broadcasts and updates the local File copy.
func (fr *FileRegistry) SyncRead() {
	ip := discovery.GetPrivateIP()
	log.Print("Read Sync enabled for: ", ip)

	// is this udp conn blocking ?
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: broadcastPort})
	if err != nil {
		log.Println("Sync Read: Error connecting to broadcast: ", err)
		return
	}

	// will this buffer get overwritten after a new Sync ?
	buff := make([]byte, 1024)
	for {
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
		fr.files[file] = ip
		fr.mu.Unlock()
		log.Print("Read Sync ended: ", ip)
	}
}

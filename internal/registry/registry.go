package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"socket-share/internal/discovery"
	"sync"
	"time"
)

type File struct {
	Name        string
	Path        string
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

// NewFile reads the absolute file path and generates a File item.
//
// It then Inserts the file into Registry.
func (fr *FileRegistry) NewFile(path string) File {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Println("Registry: Error creating new file: ", err)
	}

	file := File{
		Name:        fileInfo.Name(),
		Path:        path,
		Size:        int(fileInfo.Size()),
		Uploaded_by: discovery.GetPrivateIP(),
		Date:        int(time.Now().Unix()),
	}

	fr.Insert(file, discovery.GetPrivateIP())
	return file
}

// Insert performs a thread safe append to files map.
func (fr *FileRegistry) Insert(item File, ip string) {
	fr.mu.Lock()
	fr.files[item] = ip
	fr.mu.Unlock()
	fr.SyncWrite(item)
}

// Delete removes the file from map.
func (fr *FileRegistry) Delete(item File) {
	delete(fr.files, item)
}

// SyncWrite broadcasts the newly uploaded File item.
func (fr *FileRegistry) SyncWrite(item File) {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastIP, broadcastPort))
	if err != nil {
		log.Printf("Sync Write: Unable to resolve broadcast addr: %v", err)
		return
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 0})
	if err != nil {
		log.Println("Sync Write: Error marshaling to json: ", err)
	}

	jsonF, err := json.Marshal(item)
	if err != nil {
		log.Println("Sync Write: Error marshaling to json: ", err)
	}

	if _, err = conn.WriteToUDP(jsonF, broadcastAddr); err != nil {
		log.Println("Sync Write: Error broadcasting file for sync: ", err)
	}

	log.Printf("Broadcasted file registry for %s to %s", item.Name, broadcastAddr)
	conn.Close()
}

// SyncRead listens for File update broadcasts and updates the local File copy.
func (fr *FileRegistry) SyncRead() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: broadcastPort})
	if err != nil {
		log.Println("Sync Read: Error connecting to broadcast: ", err)
		return
	}

	// will this buffer get overwritten after a new Sync ?
	buff := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buff)
		if err != nil {
			log.Println("Sync Read: Error getting file part: ", err)
		}

		var file File
		if err := json.Unmarshal(buff[:n], &file); err != nil {
			log.Println("Sync Read: Error unmarshaling to json: ", err)
			continue
		}

		log.Printf("Sync Read: Received %s from %s", file.Name, file.Uploaded_by)
		log.Print("File Registry after sync: ", fr.files)

		if file.Uploaded_by == discovery.GetPrivateIP() {
			continue
		}

		fr.mu.Lock()
		fr.files[file] = file.Uploaded_by
		fr.mu.Unlock()
	}
}

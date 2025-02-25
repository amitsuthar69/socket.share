package registry

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"socket-share/internal/discovery"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	broadcastIP   = "192.168.0.255"
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
		log.Println("Sync Write: Conn error: ", err)
	}

	msgF, err := msgpack.Marshal(item)
	if err != nil {
		log.Println("Sync Write: Error marshaling to msgpack: ", err)
	}

	if _, err = conn.WriteToUDP(msgF, broadcastAddr); err != nil {
		log.Println("Sync Write: Error broadcasting file for sync: ", err)
	}

	log.Printf("Broadcasted file registry for %s to %s", item.Name, broadcastAddr)
	conn.Close()
}

// SyncRead listens for File update broadcasts and updates the local File copy.
func (fr *FileRegistry) SyncRead(ctx context.Context) {
	log.Print("Read Sync Started on port: ", broadcastPort)
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
		if err := msgpack.Unmarshal(buff[:n], &file); err != nil {
			log.Println("Sync Read: Error unmarshaling from msgpack: ", err)
			continue
		}

		log.Printf("Sync Read: Received %s from %s", file.Name, file.Uploaded_by)
		if file.Uploaded_by == discovery.GetPrivateIP() {
			continue
		}

		fr.mu.Lock()
		if f, exists := fr.files[file]; exists {
			log.Print("Sync Read: File Already exists: ", f)
			return
		}
		fr.files[file] = file.Uploaded_by
		fr.mu.Unlock()

		log.Print("Sync Read: File Registry after sync: ", fr.files)
		runtime.EventsEmit(ctx, "fileEvent", file)
		log.Print("Sync Read: Emitted fileEvent")
	}
}

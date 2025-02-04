package share

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

// File Server will listen for incoming file download requests from client and process the request.
func StartFileServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("File Server: Unable to start on port 8080: ", err)
		return
	}
	defer listener.Close()
	log.Println("File Server: started on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("File Server: unable to accept connection: ", err)
			continue
		}
		defer conn.Close()

		// Read the filepath from connection
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			log.Print("File Server: unable to read file name: ", err)
		}
		log.Print("File Server: Requested file name: ", string(buff[:n]))

		filename := string(buff[:n])

		/*
		  The file should be opened each time a client connects.
		  This way, the latest version is sent each time.
		*/
		file, err := os.Open(filename)
		if err != nil {
			log.Print("File Server: unable to open file: ", err)
			return
		}

		go func() {
			n, err := io.Copy(conn, file)
			if err != nil {
				log.Print("File Server: unable to copy file content to connection: ", err)
			}
			log.Printf("File Server: Transfered %d bytes to client %v", n, conn.LocalAddr())
		}()
	}
}

// File Client dials a tcp connection on provided destination IP and downloads the file provided in path.
func StartFileClient(dstIP, filePath string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dstIP, 8080))
	log.Print("file client dialing to: ", fmt.Sprintf("%s:%d", dstIP, 8080))
	if err != nil {
		log.Print("File Client: Couldn't dial to file server: ", err)
		return
	}
	defer conn.Close()

	filename := filepath.Base(filePath)
	b, err := conn.Write([]byte(filePath))
	if err != nil {
		log.Print("File Client: Couldn't send file name to server: ", err)
	}
	log.Printf("File Client: Sent file name to server: %d bytes", b)

	file, err := os.Create(filename)
	if err != nil {
		log.Print("File Client: unable to create file: ", err)
	}

	n, err := io.Copy(file, conn)
	if err != nil {
		log.Print("File Client: unable to copy file content from connection: ", err)
	}

	log.Printf("File Client: copied %d bytes into file from %v", n, conn.LocalAddr())
}

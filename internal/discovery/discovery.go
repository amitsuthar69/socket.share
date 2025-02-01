package discovery

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	broadcastPort    = 9000
	responsePort     = 9001
	broadcastMessage = "client?"
	broadcastIP      = "127.255.255.255" // change to 255 in production
)

// DiscoveryModule is responsible for peer discovery. It include two entities, client and server.
//
// Server sends a message on broadcast IP and in response clients send their subnet IP addresses.
type DiscoveryModule struct {
	clientIPs []string
	mu        sync.Mutex
	stopChan  chan bool
}

func NewDiscoveryModule() *DiscoveryModule {
	return &DiscoveryModule{
		clientIPs: make([]string, 0),
		stopChan:  make(chan bool),
	}
}

// Starts the discovery module.
//
// Runs the server and client on separate goroutines.
func (dm *DiscoveryModule) Start() {
	go dm.startServer()
	go dm.startClient()
	log.Print("discovery module started...")
}

// Stops the discovery module.
func (dm *DiscoveryModule) Stop() {
	close(dm.stopChan)
	log.Print("discovery module stopped.")
}

// Returns a list of discovered client(peers).
func (dm *DiscoveryModule) GetDiscoveredClients() []string {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	return dm.clientIPs
}

// broadcastLocationRequests sends a periodic broadcast message on the broadcast address.
//
// The broadcast address is a resolved udp4 addr produced with using broadcast IP and port.
func (dm *DiscoveryModule) broadcastLocationRequests(conn *net.UDPConn) {
	broadcastAddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastIP, broadcastPort))
	log.Print("brdc: ", broadcastAddr)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-dm.stopChan:
			return
		case <-ticker.C:
			_, err := conn.WriteToUDP([]byte(broadcastMessage), broadcastAddr)
			if err != nil {
				log.Print("error requesting IP addr: ", err)
			} else {
				log.Print("sent location request")
			}
		}
	}
}

// Server creates a udp4 listener to get client's subnet IP address.
//
// It uses the broadcastLocationRequests function to ping clients.
func (dm *DiscoveryModule) startServer() {
	listener, err := net.ListenUDP("udp4", &net.UDPAddr{Port: responsePort})
	if err != nil {
		log.Fatal("error listening on port 9001: ", err)
	}
	defer listener.Close()

	log.Printf("discovery server running on port %d", responsePort)

	go dm.broadcastLocationRequests(listener)

	for {
		select {
		case <-dm.stopChan:
			return
		default:
			// wait and listen client's response
			buff := make([]byte, 1024)
			n, addr, err := listener.ReadFromUDP(buff)
			if err != nil {
				log.Print("error reading from udp 9001: ", err)
				continue
			}

			clientIP := addr.IP.String()
			message := string(buff[:n])

			log.Print("CLIENT IP, MESSAGE: ", clientIP, message)

			dm.mu.Lock()
			if !contains(dm.clientIPs, clientIP) {
				dm.clientIPs = append(dm.clientIPs, message)
				log.Printf("New client registered: %s (reported IP: %s)", clientIP, message)
			}
			dm.mu.Unlock()
			dm.printClientList()
		}
	}
}

// client responds to server's ping by sending his subnet IP address.
func (dm *DiscoveryModule) startClient() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: broadcastPort})
	if err != nil {
		log.Fatal("error starting client: ", err)
	}
	defer conn.Close()

	log.Printf("Client listening on port %d...\n", broadcastPort)

	for {
		select {
		case <-dm.stopChan:
			return
		default:
			buff := make([]byte, 1024)
			_, addr, err := conn.ReadFromUDP(buff)
			if err != nil {
				log.Print("error reading on client: ", err)
				continue
			}

			serverAddr := &net.UDPAddr{
				IP:   addr.IP,
				Port: responsePort,
			}

			privateIP := getPrivateIP()
			_, err = conn.WriteToUDP([]byte(privateIP), serverAddr)
			if err != nil {
				log.Printf("Error sending response: %v\n", err)
			} else {
				log.Printf("Client responded to %s with IP: %s\n", serverAddr.String(), privateIP)
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (dm *DiscoveryModule) printClientList() {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	fmt.Println("Registered Clients:")
	for i, ip := range dm.clientIPs {
		fmt.Printf("%d: %s\n", i+1, ip)
	}
	fmt.Println("-------------------")
}

// A helper function which returns the subnet(private) IP address.
func getPrivateIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Print("Error getting interface addresses: ", err)
		return "unknown"
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipnet.IP
		if ip.IsLoopback() {
			continue
		}
		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}

		// Any IP in 172.x.x.x range
		if ipv4[0] == 172 {
			return ipv4.String()
		}

		if ipv4[0] == 10 ||
			// (ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31) ||
			(ipv4[0] == 192 && ipv4[1] == 168) {

			// Checks RFC 1918 private addresses
			// 10.xxx.xxx.xxx  or 10/8
			// 172.16.xxx.xxx  or 172.16/12
			// 192.168.xxx.xxx or 192.168/16

			return ipv4.String()
		}

		// Link-local addresses
		if ipv4[0] == 169 && ipv4[1] == 254 {
			return ipv4.String()
		}
	}
	return "unknown"
}

package utils

import (
	"log"
	"net"
)

// GetHostIP returns the IP address of the host that this process is currently running on
// It does so by dialing a dummy connetion and grabbing the local address from the connection,
// which indicates the public IP that the machine thinks it is communicating out from
func GetHostIP() string {
	// Make a dummy UDP connection to Google's Public DNS IP
	// Even if that IP is not addressable from this machine,
	// since it's UDP the request does not need to be valid.
	dummyConn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer dummyConn.Close()

	// Use the resolved local address's IP as the machine's preferred public IP addr
	return dummyConn.LocalAddr().(*net.UDPAddr).IP.String()
}

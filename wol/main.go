package main

import (
    "bytes"
    "fmt"
    "net"
    "os"
)

// sendMagicPacket sends a WoL magic packet to the specified MAC address.
func sendMagicPacket(macAddr string, bcastAddr string) error {
    // Parse MAC address
    mac, err := net.ParseMAC(macAddr)
    if err != nil {
        return fmt.Errorf("invalid MAC address: %v", err)
    }

    // Create the magic packet: 6 bytes of 0xFF followed by the MAC address
    // repeated 16 times
    packet := bytes.Repeat([]byte{0xFF}, 6)
    for i := 0; i < 16; i++ {
        packet = append(packet, mac...)
    }

    broadcastIP := net.ParseIP(bcastAddr)
    if broadcastIP == nil {
        return fmt.Errorf("invalid broadcast IP address: %s", bcastAddr)
    }

    // Broadcast the packet over UDP
    addr := &net.UDPAddr{
        IP:   broadcastIP,
        Port: 9, // default WoL port
    }
    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        return fmt.Errorf("failed to connect to UDP: %v", err)
    }
    defer conn.Close()

    // Send the magic packet
    _, err = conn.Write(packet)
    if err != nil {
        return fmt.Errorf("failed to send magic packet: %v", err)
    }

    return nil
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run wol.go <MAC Address> <broadcast addr>")
        return
    }

    macAddress := os.Args[1]
    bcastAddr := os.Args[2]

    fmt.Printf("Sending WoL packet to MAC address %s and bcast %s\n",
        macAddress, bcastAddr)
    if err := sendMagicPacket(macAddress, bcastAddr); err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Magic packet sent successfully.")
    }
}

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type ServerSettings struct {
	VersionName string `json:"versionName"`
	Protocol    string `json:"protocol"`
	Port        string `json:"port"`
	MOTD        string `json:"MOTD"`
}

func main() {
	var settings ServerSettings
	loadSettings(&settings)
	log.Println("Starting TCP server")
	listener, err := net.Listen("tcp", ":"+settings.Port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection: ", err)
			continue
		}
		go connectionHandler(conn, settings)
	}
}

func loadSettings(settings *ServerSettings) {
	data, err := os.ReadFile("settings.json")
	if err != err {
		panic(err)
	}
	json.Unmarshal(data, settings)
}

func connectionHandler(conn net.Conn, settings ServerSettings) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	connectionState := 0
	for {
		length, err := ReadVarInt(reader)
		if err != nil {
			return
		}
		packetData := make([]byte, length)
		_, err = io.ReadFull(reader, packetData)
		if err != nil {
			return
		}
		packetReader := bytes.NewReader(packetData)
		packetID, err := ReadVarInt(packetReader)
		if err != nil {
			return
		}

		switch packetID {
		case 0x00:
			if connectionState == 1 {
				statusString := fmt.Sprintf(`{"version":{"name":"%s","protocol":%s},"description":{"text":"%s"}}`, settings.VersionName, settings.Protocol, settings.MOTD)
				var response []byte
				response = append(response, MakeVarInt(0x00)...)
				response = append(response, MakeString(statusString)...)
				response = append(MakeVarInt(len(response)), response...)
				conn.Write(response)
				break
			}
			log.Println("---- Received handshake packet ----")
			protocolVersion, _ := ReadVarInt(packetReader)
			serverAddress, _ := ReadString(packetReader)
			serverPort, _ := ReadUnsignedShort(packetReader)
			nextState, _ := ReadVarInt(packetReader)
			log.Printf("P-ver: %d, Address: %s, Port: %d, Next State: %d\n", protocolVersion, serverAddress, serverPort, nextState)
			connectionState = nextState
		case 0x1:
			now := time.Now().UnixMilli()
			timestamp, _ := ReadLong(packetReader)
			var response []byte
			response = append(response, MakeVarInt(0x1)...)
			response = append(response, MakeLong(now)...)
			response = append(MakeVarInt(len(response)), response...)
			conn.Write(response)
			log.Printf("Ping with client: %dms\n", now-timestamp)
		default:
			log.Printf("Packet received! Length: %d, Type: %d\n", length, packetID)
		}
	}
}

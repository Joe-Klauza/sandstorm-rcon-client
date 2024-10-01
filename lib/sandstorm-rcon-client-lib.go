package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	config "github.com/Joe-Klauza/sandstorm-rcon-client/config"
	"github.com/chzyer/readline"
)

type RconPacket struct {
	Size    int32
	ID      int32
	Type    int32
	Payload string
}

var (
	conf             = config.Get()
	l         Logger = noOpLogger{}
	idCounter int32  = 0
)

func generateID() int32 {
	idCounter++
	return idCounter
}

func SetLogger(logger Logger) {
	l = logger
}

func Auth(conn net.Conn, password string) bool {
	authID := generateID()
	// Send the authentication packet.
	authPacket := BuildPacket(authID, 3, password)
	_, err := conn.Write(authPacket)
	if err != nil {
		l.Errorf("Error sending authentication packet: %s", err.Error())
		return false
	}
	// Wait for the response.
	responsePacket := ReadPacket(conn)
	if responsePacket == nil {
		l.Errorf("Error reading authentication response.")
		return false
	}
	if responsePacket.ID == authID && responsePacket.Type == 2 {
		l.Infof("Authentication successful.")
		return true
	} else {
		l.Errorf("Authentication failed.")
		return false
	}
}

func Repl(conn net.Conn) {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("help"),
		readline.PcItem("listplayers"),
		readline.PcItem("kick"),
		readline.PcItem("permban"),
		readline.PcItem("travel"),
		readline.PcItem("ban"),
		readline.PcItem("banid"),
		readline.PcItem("listbans"),
		readline.PcItem("unban"),
		readline.PcItem("say"),
		readline.PcItem("restartround"),
		readline.PcItem("maps"),
		readline.PcItem("scenarios"),
		readline.PcItem("travelscenario"),
		readline.PcItem("gamemodeproperty"),
		readline.PcItem("listgamemodeproperties"),
	)
	rl, err := readline.NewEx(&readline.Config{
		// HistoryFile:       "sandstorm-rcon-client-readline.tmp",
		Prompt:            fmt.Sprintf("RCON %s] ", conn.RemoteAddr().String()),
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		AutoComplete:      completer,
	})
	if err != nil {
		l.Errorf("Failed to instantiate readline")
		return
	}
	for {
		input, err := rl.Readline()
		if err != nil {
			break
		}
		if input == "" {
			continue
		}
		SendAndPrint(conn, input)
	}
}

func SendAndPrint(conn net.Conn, command string) error {
	output, err := Send(conn, command)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}
	// Print output
	if strings.TrimSpace(output) == "" {
		l.Infof("Server response empty")
		return nil
	}
	l.Infof("Server response:")
	fmt.Printf("%s\n", output)
	return nil
}

func Send(conn net.Conn, command string) (string, error) {
	// Generate unique IDs for the command and empty packets
	commandID := generateID()

	// Send the command packet
	commandPacket := BuildPacket(commandID, 2, command)
	l.Debugf("Sending command: %s", command)
	_, err := conn.Write(commandPacket)
	if err != nil {
		return "", fmt.Errorf("error sending command packet: %s", err.Error())
	}

	// Read and assemble response packets
	var (
		fullPayload     strings.Builder
		sentEmptyPacket bool
	)
	for {
		responsePacket := ReadPacket(conn)
		if responsePacket == nil {
			return "", fmt.Errorf("error reading command response")
		}

		if commandID != responsePacket.ID {
			l.Errorf("Received packet with unexpected ID: %d", responsePacket.ID)
			continue
		}

		if responsePacket.ID == commandID && responsePacket.Type == 0 {
			// Check if we've received the empty response packet
			if sentEmptyPacket && responsePacket.Payload == "" {
				break
			}
			fullPayload.WriteString(responsePacket.Payload)
			if !sentEmptyPacket {
				// Send the empty SERVERDATA_RESPONSE_VALUE packet
				emptyPacket := BuildPacket(commandID, 0, "")
				l.Debugf("Sending empty packet to confirm response fully received")
				_, err = conn.Write(emptyPacket)
				if err != nil {
					return fullPayload.String(), fmt.Errorf("error sending empty packet: %s", err.Error())
				}
				sentEmptyPacket = true
			}
		} else {
			return fullPayload.String(), fmt.Errorf("unhandled case in Send")
		}
	}

	return fullPayload.String(), nil
}

func BuildPacket(id int32, packetType int32, payload string) []byte {
	l.Debugf("Building packet with id %d, type %d, payload %s", id, packetType, payload)
	payloadBytes := []byte(payload)
	payloadBytes = append(payloadBytes, 0x00)      // Null terminator for the payload
	packetSize := int32(4 + 4 + len(payloadBytes)) // ID + Type + Payload
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, packetSize)
	binary.Write(buffer, binary.LittleEndian, id)
	binary.Write(buffer, binary.LittleEndian, packetType)
	buffer.Write(payloadBytes)
	return buffer.Bytes()
}

func ReadPacket(conn net.Conn) *RconPacket {
	conn.SetReadDeadline(time.Now().Add(conf.Timeout))

	// Read the packet size
	sizeBytes := make([]byte, 4)
	_, err := io.ReadFull(conn, sizeBytes)
	if err != nil {
		l.Errorf("Error reading packet size: %s", err.Error())
		return nil
	}
	packetSize := int32(binary.LittleEndian.Uint32(sizeBytes))

	// Read the rest of the packet
	packetBytes := make([]byte, packetSize)
	_, err = io.ReadFull(conn, packetBytes)
	if err != nil {
		l.Errorf("Error reading packet data: %s", err.Error())
		return nil
	}

	// Parse the packet
	packetID := int32(binary.LittleEndian.Uint32(packetBytes[0:4]))
	packetType := int32(binary.LittleEndian.Uint32(packetBytes[4:8]))
	payload := string(packetBytes[8 : len(packetBytes)-2]) // Exclude the two null terminators

	packet := &RconPacket{
		Size:    packetSize,
		ID:      packetID,
		Type:    packetType,
		Payload: payload,
	}
	l.Debugf("Received packet: ID=%d, Type=%d, Payload=%s", packetID, packetType, payload)
	return packet
}

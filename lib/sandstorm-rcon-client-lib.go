package lib

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	conf        = config.Get()
	l    Logger = noOpLogger{}
)

func SetLogger(logger Logger) {
	l = logger
}

func Auth(conn net.Conn, password string) bool {
	// Send the authentication packet.
	authPacket := BuildPacket(1, 3, password)
	_, err := conn.Write(authPacket)
	if err != nil {
		l.Errorf("Error sending authentication packet:", err.Error())
		return false
	}
	// Wait for the response.
	responsePacket := ReadPacket(conn)
	if responsePacket == nil {
		l.Errorf("Error reading authentication response.")
		return false
	}
	if responsePacket.ID != -1 {
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
		l.Infof("Server response empty", output)
		return nil
	}
	l.Infof("Server response:\n%s", output)
	return nil
}

func Send(conn net.Conn, command string) (string, error) {
	// Send the command packet.
	commandPacket := BuildPacket(1, 2, command)
	l.Debugf("Sending command: %s", command)
	_, err := conn.Write(commandPacket)
	if err != nil {
		return "", fmt.Errorf("error sending command packet: %s", err.Error())
	}
	responsePacket := ReadPacket(conn)
	if responsePacket == nil {
		return "", fmt.Errorf("error reading command response")
	}
	return responsePacket.Payload, nil
}

func BuildPacket(id int32, packetType int32, payload string) []byte {
	l.Debugf("Building packet with id %d, type %d, payload %s", id, packetType, payload)
	packetSize := int32(10 + len(payload))
	payload += "\x00\x00"
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, packetSize)
	binary.Write(buffer, binary.LittleEndian, id)
	binary.Write(buffer, binary.LittleEndian, packetType)
	buffer.WriteString(payload)
	return buffer.Bytes()
}

func ReadPacket(conn net.Conn) *RconPacket {
	conn.SetReadDeadline(time.Now().Add(conf.Timeout))
	headerBytes := make([]byte, 4*3)
	l.Debugf("Reading response header")
	_, err := conn.Read(headerBytes)
	if err != nil {
		l.Errorf("Error reading response header:", err.Error())
		return nil
	}
	packetSize := binary.LittleEndian.Uint32(headerBytes[0:4])
	packetID := binary.LittleEndian.Uint32(headerBytes[4:8])
	packetType := binary.LittleEndian.Uint32(headerBytes[8:12])
	packet := &RconPacket{
		Size:    int32(packetSize),
		ID:      int32(packetID),
		Type:    int32(packetType),
		Payload: "",
	}
	l.Debugf("Got packet with size %d, id %d, type %d", packetSize, packetID, packetType)
	if packetSize > 10 {
		payloadBytesSize := packetSize - 10
		payloadBytes := make([]byte, payloadBytesSize)
		l.Debugf("Reading response payload bytes (%d)", payloadBytesSize)
		_, err = conn.Read(payloadBytes)
		if err != nil {
			l.Errorf("Error reading response payload:", err.Error())
			return nil
		}
		packet.Payload = string(payloadBytes)
		l.Debugf("Response payload: %s", packet.Payload)
	}
	// Consume 2 bytes of padding.
	payloadBytes := make([]byte, 2)
	_, err = conn.Read(payloadBytes)
	if err != nil {
		l.Errorf("Error trimming response padding:", err.Error())
		return nil
	}
	return packet
}

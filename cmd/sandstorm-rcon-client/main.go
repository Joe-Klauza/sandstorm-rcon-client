package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/Joe-Klauza/sandstorm-rcon-client/config"
	client "github.com/Joe-Klauza/sandstorm-rcon-client/lib"
)

var conf = config.Get()

func main() {
	// Define command line flags
	server := flag.String("server", "127.0.0.1:27015", "Server IP:Port")
	password := flag.String("password", "", "RCON password")
	command := flag.String("command", "", "RCON command to execute. Omit this for REPL mode.")
	version := flag.Bool("version", false, "Print version and exit")
	// Short form flags
	flag.StringVar(server, "s", "127.0.0.1:27015", "Server IP:Port (short)")
	flag.StringVar(password, "p", "", "RCON password (short)")
	flag.StringVar(command, "c", "", "RCON command to execute (short)")
	flag.BoolVar(version, "v", false, "Print version and exit (short)")
	flag.Parse()

	if *version {
		client.Version()
		return
	}

	// Connect to the server.
	dialer := net.Dialer{Timeout: conf.Timeout}
	conn, err := dialer.Dial("tcp", *server)
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		return
	}
	fmt.Println("Connected to server:", conn.RemoteAddr().String())

	if !client.Auth(conn, *password) {
		return
	}

	if *command != "" {
		client.SendAndPrint(conn, *command)
	} else {
		client.Repl(conn)
	}

	// Close the connection.
	conn.Close()
}

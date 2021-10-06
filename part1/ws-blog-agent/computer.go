package main

import (
	"net"
	"os"
	"os/user"
)

func getDomainAndUserName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "Agent-ErrorGettingUser", err
	}
	if currentUser == nil {
		return "Agent-ErrorNilUser", err
	}
	return currentUser.Username, nil
}

func getHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "Agent-ErrorGettingHostName", err
	}
	return hostName, nil
}

func getLocalIp() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
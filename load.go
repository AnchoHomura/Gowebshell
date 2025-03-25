package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func loadConnections() {
	file, err := ioutil.ReadFile("save.txt")
	if err != nil {
		log.Printf("Failed to read save.txt: %v", err)
		return
	}
	err = json.Unmarshal(file, &connections)
	if err != nil {
		log.Printf("Failed to parse save.txt: %v", err)
	}
}

func saveConnections() {
	data, err := json.Marshal(connections)
	if err != nil {
		log.Printf("Failed to marshal connections: %v", err)
		return
	}
	err = ioutil.WriteFile("save.txt", data, 0644)
	if err != nil {
		log.Printf("Failed to write save.txt: %v", err)
	}
}

func addConnection(url, password, connType string) {
	os := ""
	var err error

	if connType == "PHP" {
		os, err = detectOs(url, password)
	} else {
		os, err = detectOs2(url, password)
	}

	if connType == "ASP" || connType == "ASPX" {
		os = "Windows"
	}

	if err != nil {
		fmt.Printf("Failed to detect OS: %v\n", err)
	}

	for _, conn := range connections {
		if conn.URL == url && conn.Password == password && conn.Type == connType {
			fmt.Println("Connection already exists.")
			return
		}
	}

	connections = append(connections, ConnectionInfo{
		URL:      url,
		Password: password,
		Type:     connType,
		OS:       os,
	})
	fmt.Println("1")
	fmt.Println(os)
	saveConnections()
}

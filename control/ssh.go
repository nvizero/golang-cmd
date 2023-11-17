package control

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ServerStep(server_date, host string) {
	statusChan <- "ServerStep1.."

	keyPath := os.ExpandEnv("/Users/tsengyenchi/.ssh/id_rsa")
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}

	statusChan <- "ServerStep2.."
	// ssh config
	config := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
		User:            "root",
	}

	// connect to ssh server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		log.Fatal("??11", err)
	}

	defer conn.Close()
	// run the ls command
	// create a new session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("??221", err)
	}
	statusChan <- "ServerStep2.5."
	defer session.Close()

	statusChan <- "ServerStep3.."
	if err := session.Run("cd /data/gameserver/ && /bin/bash start.sh"); err != nil {
		log.Fatal("??33/", err)
	}
	statusChan <- "ServerStep4.."
}

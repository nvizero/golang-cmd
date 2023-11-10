package control

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func Ma() {
	keyPath := os.ExpandEnv("/Users/tsengyenchi/.ssh/id_rsa")
	keyBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read private key: %s", err)
	}

	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %s", err)
	}

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
	conn, err := ssh.Dial("tcp", "173.249.199.139:22", config)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	// run the ls command
	// create a new session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// run the ls command
	output, err := session.CombinedOutput("ls")
	if err != nil {
		log.Fatal(err)
	}

	// print the output of ls command
	log.Printf("ls command output:\n%s", output)
}

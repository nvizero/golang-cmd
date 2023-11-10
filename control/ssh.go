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

	if err := session.Run("cd /home"); err != nil {
		log.Fatal(err)
	}
	// create a new session for ls command
	session, err = conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// list contents of the directory
	output, err := session.CombinedOutput("ls")
	if err != nil {
		log.Fatal(err)
	}
	statusChan <- string(output)
	log.Printf("Contents of /home:\n%s", output)

	// create a new session for php -v command
	session, err = conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// run the php -v command
	output, err = session.CombinedOutput("pwd")
	if err != nil {
		log.Fatal(err)
	}
	statusChan <- string(output)
	log.Printf("php -v command output:\n%s", output)
}

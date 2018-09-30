package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	cmd := os.Args[1]
	hosts := os.Args[2:]

	results := make(chan string, 10)
	timeout := time.After(10 * time.Second)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "22"
	}

	config := &ssh.ClientConfig{
		// User: os.Getenv("USER"),
		User: "root",
		Auth: []ssh.ClientAuth{makeKeyring()},
	}

	for _, hostname := range hosts {
		go func(hostname string, port string) {
			results <- executeCmd(cmd, hostname, port, config)
		}(hostname, port)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}

type SignerContainer struct {
	signers []ssh.Signer
}

func (t *SignerContainer) Key(i int) (key ssh.PublicKey, err error) {
	if i >= len(t.signers) {
		return
	}
	key = t.signers[i].PublicKey()
	return
}

func (t *SignerContainer) Sign(i int, rand io.Reader, data []byte) (sig []byte, err error) {
	if i >= len(t.signers) {
		return
	}
	sig, err = t.signers[i].Sign(rand, data)
	return
}

func makeSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func makeKeyring() ssh.ClientAuth {
	signers := []ssh.Signer{}
	// keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa"}
	keys := []string{os.Getenv("HOME") + "/Downloads/personals/shammi.pem"}

	for _, keyname := range keys {
		signer, err := makeSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.ClientAuthKeyring(&SignerContainer{signers})
}

func executeCmd(command, hostname string, port string, config *ssh.ClientConfig) string {
	conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return fmt.Sprintf("%s -> %s", hostname, stdoutBuf.String())
}

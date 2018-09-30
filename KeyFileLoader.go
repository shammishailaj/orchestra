package orchestra

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

// Code to log-in via key file
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %s", err)
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		fmt.Printf("Error parsing private key: %s", err)
		return nil
	}
	return ssh.PublicKeys(key)
}

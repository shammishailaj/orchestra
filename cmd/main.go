package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"goprojects/orchestra"
	"os"
)

func main() {
	fmt.Printf("\n\n")

	userName := flag.String("u", "", "Username of the account which is to be used to log-in.")
	keyFile := flag.String("k", "", "Path to public key file.")
	password := flag.String("p", "", "Password for the user account being used.")
	hostName := flag.String("h", "", "Hostname of the server to connect to.")
	portNumber := flag.Int("port", 22, "Port number of SSH Server on the remote host.")
	authMethod := flag.String("am", "key", "Authentication Method to be used. \"key\" for Public Key, \"password\" for password based authentication.")
	flag.Parse()

	if *userName == "" {

		fmt.Printf("Must supply username for the remote machine \"key\"\n\n")
		flag.PrintDefaults()

	} else {

		if *hostName == "" {
			fmt.Printf("Must supply hostname of the remote server\n")
			flag.PrintDefaults()
			os.Exit(1)
		}

		// var hostKey, hostKeyErr = getHostKey(*hostName)
		// if hostKeyErr != nil {
		// 	fmt.Printf("HostKeyErr: %s\n", hostKeyErr)
		// }

		switch *authMethod {
		case "key":
			if *keyFile == "" {
				fmt.Printf("Must supply path to public key file with auth method \"key\"\n\n")
				flag.PrintDefaults()
			} else {
				// Code to log-in via key file
				sshConfig := &ssh.ClientConfig{
					User: *userName,
					Auth: []ssh.AuthMethod{
						orchestra.PublicKeyFile(*keyFile),
					},
					// HostKeyCallback: ssh.FixedHostKey(orchestra.HostKeyCheck(*hostName)),
					HostKeyCallback: ssh.InsecureIgnoreHostKey(),
				}

				orchestra.ConnectSSH(*hostName, *portNumber, sshConfig)
			}
		case "password":
			if *password == "" {
				fmt.Printf("Must supply password with auth method \"password\"\n\n")
				flag.PrintDefaults()
			} else {
				// Code for logging-in via password
				sshConfig := &ssh.ClientConfig{
					User: *userName,
					Auth: []ssh.AuthMethod{
						ssh.Password(*password),
					},
				}

				orchestra.ConnectSSH(*hostName, *portNumber, sshConfig)

			}
		default:
			fmt.Printf("Unknown auth method \"%s\"\n\n", *authMethod)
			flag.PrintDefaults()
		}
	}
	fmt.Printf("\n\n")
}

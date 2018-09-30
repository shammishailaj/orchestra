package orchestra

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

func ConnectSSH(host string, port int, sshConfig *ssh.ClientConfig) {
	fmt.Printf("Dialing SSH to %s:%d...", host, port)
	hostPort := fmt.Sprintf("%s:%d", host, port)
	connection, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		fmt.Printf("Failed to dial: %s", err)
	} else {
		fmt.Printf("CONNECTED!\n Creating SSH Session...")
		session, err := connection.NewSession()
		if err != nil {
			fmt.Printf("Failed to create session: %s", err)
		} else {
			fmt.Printf("CREATED!\n Requesting PseudoTerminal...")
			modes := ssh.TerminalModes{
				ssh.ECHO:          0,     // disable echoing
				ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
				ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
			}

			if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
				session.Close()
				fmt.Printf("request for pseudo terminal failed: %s", err)
			} else {
				fmt.Printf("RECEIVED!\n Attaching to STDIN...")
				stdin, err := session.StdinPipe()
				if err != nil {
					fmt.Printf("Unable to setup stdin for session: %v", err)
				} else {
					fmt.Println("ATTACHED!\n")
				}
				go io.Copy(stdin, os.Stdin)

				fmt.Printf("Attaching to STDOUT...")
				stdout, err := session.StdoutPipe()
				if err != nil {
					fmt.Printf("Unable to setup stdout for session: %v", err)
				} else {
					fmt.Println("ATTACHED!\n")
					fmt.Println("Stdout for current session has been piped")
				}
				go io.Copy(os.Stdout, stdout)

				fmt.Printf("Attaching to STDERR...")
				stderr, err := session.StderrPipe()
				if err != nil {
					fmt.Printf("Unable to setup stderr for session: %v", err)
				} else {
					fmt.Println("ATTACHED!\n")
					fmt.Println("Stderr for current session has been piped")

				}
				go io.Copy(os.Stderr, stderr)
			}
		}
	}
}

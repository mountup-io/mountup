/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh <clientname>",
	Short: "ssh into your virtual machine",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a clientname")
		}
		if len(args) > 1 {
			return errors.New("too many arguments, only one clientname is required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Should check if table exists, and it there's an error

		db := api.NewDB()
		host, err := db.GetHostnameForClientname(args[0])
		if err != nil {
			//Either table or client with this name doesn't exist
			fmt.Printf("virtual machine named %s not found.\nTo create a client run:\n\tmountup create %s\n", args[0], args[0])
			return
		}

		homeDir, _ := os.UserHomeDir()
		pkeyDir := filepath.Join(homeDir, ".mountup/keys", args[0]+".pem")
		publicKey, err := PublicKeyFile(pkeyDir)
		if err != nil {
			log.Println(err)
			return
		}

		config := &ssh.ClientConfig{
			User: "ubuntu",
			Auth: []ssh.AuthMethod{
				publicKey,
			},
			Timeout:         5 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		conn, err := ssh.Dial("tcp", host+":22", config)
		if err != nil {
			panic("Failed to dial: " + err.Error())
		}
		defer conn.Close()

		session, err := conn.NewSession()
		if err != nil {
			panic("Failed to create session: " + err.Error())
		}
		defer session.Close()

		// Set IO
		session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
		session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
		session.Stdin = os.Stdin

		modes := ssh.TerminalModes{
			ssh.ECHO:  0, // Disable echoing
			ssh.IGNCR: 1, // Ignore CR on input.
		}

		// Request pseudo terminal
		if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
			log.Fatalf("request for pseudo terminal failed: %s", err)
		}

		// Start remote shell
		if err := session.Shell(); err != nil {
			log.Fatalf("failed to start shell: %s", err)
		}

		// Accepting commands
		if err := session.Wait(); err != nil {
			if e, ok := err.(*ssh.ExitError); ok {
				switch e.ExitStatus() {
				case 130:
					fmt.Print("Session ended")
					return
				}
			}
			fmt.Errorf("ssh: %s", err)
			return
		}
		fmt.Printf("Connection to %s closed.\n", host)
		return
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

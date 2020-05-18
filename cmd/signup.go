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
	"bufio"
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

// signupCmd represents the signup command
var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up a new account",
	Long:  `Sign up a new account`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')

		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("failed to get password")
		}
		password := string(bytePassword)

		err = api.MakeSignUpRequest(username, email, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("\nThanks for signing up!")
	},
}

func init() {
	// DISABLED
	//rootCmd.AddCommand(signupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

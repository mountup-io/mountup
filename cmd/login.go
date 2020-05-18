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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mountup-io/mountup/api"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"os"
	"syscall"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your mountup account",
	Long:  `Login to your mountup account`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		username, _ := reader.ReadString('\n')

		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("failed to get password")
		}
		password := string(bytePassword)

		resp, err := makeLoginRequest(username, password)
		if err != nil {
			fmt.Println("error making login request")
			return
		}

		db := api.NewDB()

		err = db.InitTables()
		if err != nil {
			fmt.Printf("Error initializing mountup db: %s\n", err)
			return
		}

		// Insert access and refresh tokens
		tokenDetails := api.TokenDetails{}
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "access_token" {
				tokenDetails.AccessToken = cookie.Value
				tokenDetails.AccessTokenExp = cookie.Expires
			}
			if cookie.Name == "refresh_token" {
				tokenDetails.RefreshToken = cookie.Value
				tokenDetails.RefreshTokenExp = cookie.Expires
			}
		}
		err = db.PutAuthTokens(tokenDetails)
		if err != nil {
			fmt.Printf("Error saving auth tokens into mountup local db: %s\n", err)
			return
		}

		fmt.Printf("Logged in as %s\n", username)
	},
}

func init() {
	// DISABLED
	//rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func makeLoginRequest(username string, password string) (*http.Response, error) {
	url := "http://localhost:8080/login"

	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return resp, nil
}

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
	"strings"
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
		username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		bytePassword, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			fmt.Println("failed to get password")
		}
		password := string(bytePassword)
		fmt.Printf("\n")

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

		err = db.DeleteAuthTokens()
		if err != nil {
			return
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
	rootCmd.AddCommand(loginCmd)
}

func makeLoginRequest(username string, password string) (*http.Response, error) {
	url := "http://api.mountup.io/login"

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
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

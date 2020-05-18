package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/machinebox/graphql"
	"log"
	"net/http"
	"time"
)

type CreateVMForUserAccount struct {
	Vm   VM         `json:"vm"`
	Pkey PrivateKey `json:"pkey"`
}

type CreateVMForUserAccountResponseStruct struct {
	CreateVMForUserAccount CreateVMForUserAccount `json:"createVMForUserAccount"`
}

func MakeCreateVMRequest(clientName string) (*VM, *PrivateKey, error) {
	client := graphql.NewClient("http://localhost:8080/graphql")

	// make a request
	req := graphql.NewRequest(`
			mutation createVMForUserAccount($clientName: String!) {
			  createVMForUserAccount(input: {clientName: $clientName}){
				vm {
				  ID
				  ownerID
				  name
				  instanceID
				  publicDNS
				  publicIP
				  keyName
				  zone
				}
				pkey {
				  keyName
				  keyMaterial
				  keyFingerprint
				}
			  }
			}
		`)

	// set any variables
	req.Var("clientName", clientName)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	authToken, err := getAuthToken()
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData CreateVMForUserAccountResponseStruct
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	return &respData.CreateVMForUserAccount.Vm, &respData.CreateVMForUserAccount.Pkey, nil
}

func getAuthToken() (string, error) {
	db := NewDB()

	accessToken, refreshToken, accessTokenExp, refreshTokenExp, err := db.GetAuthTokens()

	// Check validity of access_token
	atExp, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", accessTokenExp)
	if err != nil {
		fmt.Println("failed to parse token expiration date")
		return "", err
	}
	if atExp.After(time.Now()) {
		return accessToken, nil
	}

	// Check validity of refresh_token
	rtExp, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", refreshTokenExp)
	if err != nil {
		fmt.Println("failed to parse token expiration date")
		return "", err
	}

	// Query /refresh and get new tokens
	if rtExp.After(time.Now()) {

		resp, err := makeRefreshRequest(refreshToken)
		if err != nil {
			fmt.Println("error making login request")
			return "", err
		}

		// Insert access and refresh tokens
		tokenDetails := TokenDetails{}
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
			return "", err
		}

		err = db.PutAuthTokens(tokenDetails)
		if err != nil {
			return "", err
		}

		return tokenDetails.AccessToken, nil
	}

	// Otherwise fail and tell user to relogin
	return "", errors.New("auth tokens expired, please relogin")
}

func makeRefreshRequest(refreshToken string) (*http.Response, error) {
	url := "http://localhost:8080/refresh"

	reqBody, err := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
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

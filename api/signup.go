package api

import (
	"context"
	"github.com/machinebox/graphql"
	"log"
)

type SignUpUserAccount struct {
	Token string `json:"token"`
}

type MakeSignUpRequestStruct struct {
	SignUpUserAccount SignUpUserAccount `json:"signUpUserAccount"`
}

func MakeSignUpRequest(username string, email string, password string) error {
	client := graphql.NewClient("http://api.mountup.io:8080/graphql")

	// make a request
	req := graphql.NewRequest(`
			mutation SignUpUserAccount($username: String!, $password: String!, $email: String!) {
			  signUpUserAccount(input: {username: $username, password: $password, email: $email}){
				token
			  }
			}
		`)

	// set any variables
	req.Var("username", username)
	req.Var("email", email)
	req.Var("password", password)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData MakeSignUpRequestStruct
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	return nil
}

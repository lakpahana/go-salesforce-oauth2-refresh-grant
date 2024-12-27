package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/simpleforce/simpleforce"
	"golang.org/x/oauth2"
)

func loadEnvVars() (clientID, clientSecret, refreshToken, tokenURL, instanceURL string, err error) {
	clientID = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	refreshToken = os.Getenv("REFRESH_TOKEN")
	tokenURL = os.Getenv("TOKEN_URL")
	instanceURL = os.Getenv("INSTANCE_URL")

	if clientID == "" || clientSecret == "" || refreshToken == "" || tokenURL == "" || instanceURL == "" {
		err = fmt.Errorf("missing one or more required environment variables")
	}

	return
}

func getAccessTokenUsingRefreshToken(clientID, clientSecret, refreshToken, tokenURL string) (string, error) {

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
		},
	}

	tok := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := conf.TokenSource(context.Background(), tok)
	newToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %v", err)
	}

	return newToken.AccessToken, nil
}

func querySalesforce(client *simpleforce.Client) ([]simpleforce.SObject, error) {

	query := "SELECT Id, Account.Name, Name FROM Contact LIMIT 5"

	records, err := client.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query Salesforce: %v", err)
	}

	return records.Records, nil
}

func main() {

	clientID, clientSecret, refreshToken, tokenURL, instanceURL, err := loadEnvVars()
	if err != nil {
		log.Fatal(err)
	}

	accessToken, err := getAccessTokenUsingRefreshToken(clientID, clientSecret, refreshToken, tokenURL)
	if err != nil {
		log.Fatal(err)
	}

	client := simpleforce.NewClient(instanceURL, clientID, simpleforce.DefaultAPIVersion)
	client.SetSidLoc(accessToken, instanceURL)

	records, err := querySalesforce(client)
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		fmt.Println(record)
	}
}

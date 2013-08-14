// Based on goauth2 example :
//   Copyright 2011 The goauth2 Authors. All rights reserved.
//   Use of this source code is governed by a BSD-style
//   license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	clientId     = flag.String("id", "your client id", "Client ID")
	clientSecret = flag.String("secret", "your client secret", "Client Secret")
	redirectURL  = flag.String("redirect_url", "http://example.net", "Redirect URL")
	authURL      = flag.String("auth_url", "https://na.airvantage.net/api/oauth/authorize", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://na.airvantage.net/api/oauth/token", "Token URL")
	requestURL   = flag.String("request_url", "https://na.airvantage.net/api/gateways", "API request")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, create a new client in Develop=>API client menu.

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

// Some simple structures representing the result of API calls

// a Gateway
type Gateway struct {
	State        string
	Type         string
	Metadata     []KeyValue
	Labels       []string
	CreationDate int64
	Imei         string
	SerialNumber string
	MacAddress   string
	Uid          string
}

// key value struct for metadata
type KeyValue struct {
	Key   string
	Value string
}

// the result of the gateway API call
type GatewayResult struct {
	Items  []Gateway
	Count  int32
	Size   int32
	Offset int32
}

// to String for the Gateway struct
func (gw Gateway) String() string {
	var creationDate = time.Unix(gw.CreationDate/1000, 0)
	return fmt.Sprintf("uid: %s state: %s type:%s imei: %s s/n: %s creation: %s ", gw.Uid, gw.State, gw.Type, gw.Imei, gw.SerialNumber, creationDate.String())
}

// to String for the GatewayResult struct
func (gr GatewayResult) String() string {
	var res = fmt.Sprintf("GatewayResult:\n   count: %d\n   size: %d\n   offset: %d\n   gateways:\n", gr.Count, gr.Size, gr.Offset)

	for _, gw := range gr.Items {
		res += "      " + gw.String() + "\n"
	}

	res += "\n"
	return res
}

// main program entry point
func main() {
	fmt.Println("hello AirVantage, will try to authenticate and get a list of the present gateways")
	flag.Parse()

	// Set up a OAuth configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			return
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token

	// Make the request.
	r, err := transport.Client().Get(*requestURL)
	if err != nil {
		log.Fatal("Get:", err)
	}
	// be sure to close the stream
	defer r.Body.Close()

	// create a JSon decoder for the HTTP request result
	dec := json.NewDecoder(r.Body)
	var gwResult GatewayResult
	if err := dec.Decode(&gwResult); err != nil {
		log.Fatal("JSon decode:", err)
	}
	fmt.Println("result : ")
	fmt.Println(gwResult)
}

airvantage-api-golang
=====================

AirVantage Go example

A simple OAuth2 + JSon client for gathering the list of gateway from http://airvantage.net

Based on https://code.google.com/p/goauth2/source/browse/oauth/example/oauthreq.go

Building
--------

First create a new client in the AirVantage "Develop : API Client" screen.
And complete main.go with the generated client id and secret.

Install the goauth2 library
> go get code.google.com/p/goauth2/oauth

Build
> go build

An executable should be present in the current directory, run it! 


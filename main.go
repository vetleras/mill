package main

import (
	log "github.com/sirupsen/logrus"
	mill "github.com/vetleras/mill/mill"
)

func init(){
	log.SetLevel(log.DebugLevel)
}

func main() {
	filepath := ".mill.json"
	credentials, err := mill.CredentialsFromFile(filepath)

	if err != nil {
		log.Error(err)
		credentials, err = mill.CredentialsFromPrompt()
		if err != nil {
			log.Fatal(err)
		}
	}

	//do stuff

	err = credentials.ToFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
}

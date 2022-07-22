package main

import (
	"log"
)

func main() {
	err := Load(".env")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = Validate()
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := DB{}
	err = db.connectToDB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.client.Close()
}

package main

import (
	"log"
)

func main() {
	err := loadENV()
	if err != nil {
		log.Fatalf(err.Error())
	}

	client, err := connectToDB()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.Close()
}

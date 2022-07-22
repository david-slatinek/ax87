package main

import (
	"log"
	"os"
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

	db := DB{org: os.Getenv("ORGANIZATION"), bucket: os.Getenv("BUCKET")}
	err = db.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.client.Close()

	db.Add(&Data{
		DataType: carbonMonoxide,
		Value:    10.78,
	})
}

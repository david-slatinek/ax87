package main

import (
	"fmt"
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

	//err = db.Init()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//db.Add(&Data{
	//	DataType: carbonMonoxide,
	//	Value:    450.78,
	//})	res, err := db.Last24H(carbonMonoxide)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println()
	//	fmt.Println(res)
	//
	//db.Add(&Data{
	//	DataType: carbonMonoxide,
	//	Value:    102.1,
	//})

	//res, err := db.Latest(carbonMonoxide)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	res, err := db.Last24H(carbonMonoxide)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println(res)
}

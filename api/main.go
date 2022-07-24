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

	db := DB{url: os.Getenv("INFLUXDB_URL"), token: os.Getenv("INFLUXDB_TOKEN"),
		org: os.Getenv("ORGANIZATION"), bucket: os.Getenv("BUCKET")}
	err = db.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.client.Close()

	//err = db.Init()
	//if err != nil {
	//	fmt.Println(err)
	//}

	//db.Add(&Data{
	//	DataType: carbonMonoxide,
	//	Value:    40.78,
	//})
	//
	//db.Add(&Data{
	//	DataType: airQuality,
	//	Value:    120.4,
	//})
	//
	//db.Add(&Data{
	//	DataType: raindrops,
	//	Value:    200,
	//})
	//
	//db.Add(&Data{
	//	DataType: soilMoisture,
	//	Value:    300,
	//})

	//res, err := db.Latest(soilMoisture)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//res, err := db.Last24H(soilMoisture)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//res, err := db.Median(soilMoisture)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//res, err := db.Max(soilMoisture)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)

	//res, err := db.Min(soilMoisture)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(res)
}

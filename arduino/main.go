package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
)

func readFromSerial(s *serial.Port, id int) (float32, error) {
	n, err := s.Write([]byte{byte(id)})
	fmt.Println("First n:", n)
	if err != nil {
		return 0, err
	}

	//err = s.Flush()
	//if err != nil {
	//	log.Println(err)
	//}

	reader := bufio.NewReader(s)
	result, err := reader.ReadBytes('\x00')
	if err != nil {
		return 0, err
		//log.Println(err)
	}
	fmt.Println(result)
	//return result, nil

	//fmt.Println("String:", string(buf[:]))
	return 0, nil
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	readPtr := flag.Int("read", 0, "Read sensor value/values\n\t0 = all (default)\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")
	//printPtr := flag.Bool("print", true, "Print read value/values")

	flag.Parse()

	if *readPtr < 0 || *readPtr > 4 {
		log.Println("Invalid value for -read argument, exiting")
		os.Exit(1)
	}

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("serial.OpenPort")
		log.Fatal(err)
	}

	defer func(s *serial.Port) {
		handleError(s.Close())
	}(s)

	//fmt.Println(*printPtr)
}

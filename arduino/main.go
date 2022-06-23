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

func readFromSerial(s *serial.Port, code int) (float32, error) {
	n, err := s.Write([]byte("1"))
	fmt.Println("First n:", n)
	if err != nil {
		fmt.Println("s.Write([]byte(\"1\"))")
		log.Fatal(err)
	}

	err = s.Flush()
	if err != nil {
		log.Println(err)
	}

	reader := bufio.NewReader(s)
	reply, err := reader.ReadBytes('\x00')
	if err != nil {
		log.Println(err)
	}
	fmt.Println(reply)

	//fmt.Println("String:", string(buf[:]))
	return 0, nil
}

func readCarbonMonoxide(s *serial.Port) (float32, error) {
	return readFromSerial(s, 1)
}

func main() {
	readPtr := flag.Int("read", 0, "Read sensor value/values\n\t0 = all (default)\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")
	//printPtr := flag.Bool("print", true, "Print read value/values")

	flag.Parse()

	if *readPtr < 0 || *readPtr > 4 {
		log.Println("Invalid value for -read argument!")
		log.Println("Exiting...")
		os.Exit(1)
	}

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("serial.OpenPort")
		log.Fatal(err)
	}
	defer func(s *serial.Port) {
		err := s.Close()
		if err != nil {
			fmt.Println("s.Close()")
			log.Println(err)
		}
	}(s)

	_, err = readCarbonMonoxide(s)
	if err != nil {
		fmt.Println("readCarbonMonoxide(s)")
		log.Println(err)
	}

	//fmt.Println(*printPtr)
}

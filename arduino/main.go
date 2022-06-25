package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func readFromSerial(s *serial.Port, id int) (float32, error) {
	_, err := s.Write([]byte(strconv.Itoa(id)))
	if err != nil {
		return 0, err
	}

	reader := bufio.NewReader(s)
	result, err := reader.ReadBytes('\x00')
	if err != nil {
		return 0, err
	}

	res := strings.TrimSpace(string(result))
	res = strings.TrimSuffix(res, "\x00")

	nr, err := strconv.ParseFloat(res, 32)
	if err != nil {
		return 0, err
	}

	return float32(nr), nil
}

func main() {
	readPtr := flag.Int("read", 0, "Read sensor value/values\n\t0 = all (default)\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")

	flag.Parse()

	if *readPtr < 0 || *readPtr > 4 {
		log.Println("Invalid value for -read argument, exiting")
		os.Exit(1)
	}

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer func(s *serial.Port) {
		err := s.Close()
		if err != nil {
			log.Println(err)
		}
	}(s)

	if *readPtr == 0 {
		for i := 1; i <= 4; i++ {
			data, err := readFromSerial(s, i)
			log.Println("Argument:", i)

			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("Result:", data)
			}

			time.Sleep(time.Second)
		}
		return
	}

	data, err := readFromSerial(s, *readPtr)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Result:", data)
	}
}

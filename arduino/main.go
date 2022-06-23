package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	readPtr := flag.Int("read", 0, "Read sensor value/values\n\t0 = all (default)\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")
	printPtr := flag.Bool("print", true, "Print read value/values")

	flag.Parse()

	if *readPtr < 0 || *readPtr > 4 {
		log.Println("Invalid value for -read argument!")
		log.Println("Exiting...")
		os.Exit(1)
	}

	fmt.Println(*readPtr)
	fmt.Println(*printPtr)
}

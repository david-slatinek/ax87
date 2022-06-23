package main

import (
	"flag"
	"fmt"
	"os"
)

const version float32 = 1.0

func main() {
	readPtr := flag.Int("read", 0, "Read sensor value.\n\t0 = all\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")
	printPtr := flag.Bool("print", true, "Print values")
	versionPtr := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *versionPtr {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

	fmt.Println(*readPtr)
	fmt.Println(*printPtr)
}

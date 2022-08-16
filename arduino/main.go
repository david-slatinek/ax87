package main

import (
	pb "arduino/schema"
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"github.com/tarm/serial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
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

func Upload(dataType int, value float32, client pb.RequestClient) (*pb.Reply, error) {
	return client.Add(context.Background(), &pb.Data{
		DataType:  pb.DataType(dataType),
		Value:     value,
		Timestamp: timestamppb.New(time.Now()),
	})
}

func LoadTLS() (credentials.TransportCredentials, error) {
	serverCa, err := ioutil.ReadFile("../cert/ca-cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverCa) {
		return nil, errors.New("failed to add server CA certificate")
	}

	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	//readPtr := flag.Int("read", 0, "Read sensor value/values\n\t0 = all (default)\n\t1 = carbon monoxide\n\t2 = air quality\n\t3 = raindrops\n\t4 = soil moisture")
	uploadPtr := flag.Bool("upload", false, "Upload read data to the API")

	flag.Parse()

	//if *readPtr < 0 || *readPtr > 4 {
	//	log.Fatal("Invalid value for -read argument, exiting")
	//}

	//c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600, ReadTimeout: time.Second * 5}
	//s, err := serial.OpenPort(c)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func(s *serial.Port) {
	//	if err := s.Close(); err != nil {
	//		log.Println(err)
	//	}
	//}(s)

	tlsCred, err := LoadTLS()
	if err != nil && *uploadPtr {
		log.Fatal(err)
	}

	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithTransportCredentials(tlsCred))
	if err != nil {
		log.Println(err)
		if *uploadPtr {
			os.Exit(1)
		}
	}

	defer func(conn *grpc.ClientConn) {
		if err != nil {
			return
		}

		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}(conn)

	client := pb.NewRequestClient(conn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//if *readPtr == 0 {
	//	for i := 1; i <= 4; i++ {
	//		data, err := readFromSerial(s, i)
	//		fmt.Println("Argument:", i)
	//
	//		if err != nil {
	//			log.Println(err)
	//		} else {
	//			fmt.Println("Result:", data)
	//
	//			if *uploadPtr {
	//				_, err := Upload(i, data, client)
	//				if err != nil {
	//					log.Println(err)
	//				}
	//			}
	//		}
	//		time.Sleep(time.Second)
	//	}
	//	return
	//}

	//data, err := readFromSerial(s, *readPtr)
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	fmt.Println("Result:", data)
	//
	//	if *uploadPtr {
	//		_, err := Upload(*readPtr, data, client)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//	}
	//}

	if *uploadPtr {
		_, err := Upload(1, 45.4, client)
		if err != nil {
			log.Println(err)
		}
	}
}

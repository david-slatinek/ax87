![Project logo](/images/logo.png)

# Table of Contents

- [Table of Contents](#table-of-contents)
- [About](#about)
- [Arduino](#arduino)
- [API](#api)

# About

Smart home system with Arduino, InfluxDB, Go, gRPC, Docker, and Ionic.

The project consists of 4 main components:
- Data capture from the sensors using Arduino.
- Database for data storage.
- API for storing and retrieving data.
- Mobile app for displaying the values.

# Arduino
<div align="center">
  <img alt="Arduino" src="https://img.shields.io/badge/Arduino-00979D?style=for-the-badge&logo=Arduino&logoColor=white"/>
</div>

Arduino is responsible for getting the data from the sensors for the following:
- Carbon monoxide.
- Air quality.
- Raindrops.
- Soil moisture.

After that, it uploads data to the API by calling an appropriate **grpc** method:
```go
conn, err := grpc.Dial("address", grpc.WithTransportCredentials(tlsCred))
client := pb.NewRequestClient(conn)

client.Add(context.Background(), &pb.Data{
    DataType:  pb.DataType(dataType),
    Value:     value,
    Timestamp: timestamppb.New(time.Now()),
})
```

<br>

Sensor schematic:
<div align="center">
  <img src="./images/circuit-designs/design-all.png" alt="Arduino wiring" height="500" width="700">
</div>

More circuit design images can be seen [here](/images/circuit-designs/).

# API
<div align="center">
  <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/>
  <img alt="Shell Script" src="https://img.shields.io/badge/Shell_Script-121011?style=for-the-badge&logo=gnu-bash&logoColor=white"/>
  <img alt="openssl" src="https://img.shields.io/badge/OpenSSL-721412?style=for-the-badge&logo=openssl&logoColor=white"/>
  <img alt="InfluxDB" src="https://img.shields.io/badge/InfluxDB-22ADF6?style=for-the-badge&logo=InfluxDB&logoColor=white"/>
  <img alt="redis" src="https://img.shields.io/badge/redis-CC0000.svg?&style=for-the-badge&logo=redis&logoColor=white"/>
</div>

The API was made with the language **Go** and with the **grpc** framework. For the database, we chose **influxdb**, which is hosted by the **influxdb cloud**.

Data is stored in a single bucket with four different measurements:
- CARBON_MONOXIDE.
- AIR_QUALITY.
- RAINDROPS.
- SOIL_MOISTURE.

Each measurement has three fields: *value*, i. e. the recorded value from the sensor; *_time*, i. e. when the *value* was taken; and *category*, with it being in the following range:
- For CARBON_MONOXIDE it is [1, 7], with 1 being the best.
- For AIR_QUALITY it is [1, 6], with 1 being the best.
- For RAINDROPS it is [1, 4], with 1 indicating no or little rain.
- For SOIL_MOISTURE it is [0, 100]%, with 0 indicating no soil moisture.

To speed up queries, we used **Redis** as an in-memory cache, which stores the latest added data. We also added support for unit testing, with the appropriate files having the *_test.go* suffix.

Lastly, we added support for mutual TLS with the help of **OpenSSL**.

The grpc service (*.proto* file) can be seen [here](/api/schema/).

Method to get the latest record:
```go
func (server *Server) Latest(_ context.Context, 
request *pb.DataRequest) (*pb.DataWithCategory, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request can't be nil")
	}

	var latest *model.DataResponse
	var ok = false

	if server.cache != nil {
		value, err := server.cache.Get(request.DataType.String())

		if err != nil && server.Development {
			log.Printf("Error with cache get, error: %v", err)
		} else {
			if err := json.Unmarshal([]byte(value), &latest); 
			err != nil && server.Development {
				log.Printf("Error with json.Unmarshal, error: %v", err)
			} else if err != redis.Nil && value != "" {
				ok = true
			}
		}
	}

	if !ok {
		var err error

		latest, err = server.DBService.Latest(request.GetDataType().String())
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		if err = server.AddToCache(latest); err != nil && server.Development {
			log.Printf("Error with cache set, error: %v", err)
		}
	}

	return latest.Convert(), nil
}
```

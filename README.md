![Project logo](/images/logo.png)

# Table of Contents

- [Table of Contents](#table-of-contents)
- [About](#about)
- [Arduino](#arduino)

# About

Smart home system with Arduino, InfluxDB, Go, Docker, gRPC, and Ionic.

The project consists of 4 main components:
- Data capture from the sensors using Arduino
- Database for data storage
- API for storing and retrieving data
- Mobile app for displaying the values

# Arduino
<div align="center">
  <img alt="Arduino" src="https://img.shields.io/badge/Arduino-00979D?style=for-the-badge&logo=Arduino&logoColor=white"/>
</div>

Arduino is responsible for getting the data from the sensors for the following:
- Carbon monoxide
- Air quality
- Raindrops
- Soil moisture

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

<br>

More circuit design images can be seen [here](/images/circuit-designs/).


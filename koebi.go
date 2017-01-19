package main

import (
    "flag"
    "log"
    "net/http"
    "github.com/dsbrng25b/koebi/sensor"
    "encoding/json"
    "time"
    "github.com/influxdata/influxdb/client/v2"
)

var (
    sensorPin = flag.Int("sensor.pin", 4, "pin where the sensor is connected")
    sensorMeasureInterval = flag.Duration("sensor.interval", time.Duration(3) * time.Second, "time interval between measurements")
    sensorType = flag.String("sensor.type", "DHT11", "type of the sensor. either DHT11 or DHT22")
    listenAddress = flag.String("web.addr", ":80", "Address on which to expose the http api.")
    influxDbUrl = flag.String("influxdb.url", "http://localhost:8086", "url of the influx database")
    influxDbName = flag.String("influxdb.name", "temperature", "name of the influx database")
    influxDbUser = flag.String("influxdb.user", "root", "user to connect to the influx database")
    influxDbPassword = flag.String("influxdb.password", "root", "password to connect to the influx database")
    locationTag = flag.String("tag.location", "", "string which specifies the location of the sensor (e.g. home)")
    roomTag = flag.String("tag.room","", "string which specifies the location of the sensor (e.g. home)")
)


func main() {
    flag.Parse()

    mySensor, err := sensor.New(*sensorType, *sensorPin, *sensorMeasureInterval)
    if err != nil {
        log.Fatalln("could not create sensor: ", err)
    }

    //start measurements
    go mySensor.Start()

    go dbWriter(mySensor)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        m, err := mySensor.GetLast()
        if err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "no measurement available"})
        }else{
            json.NewEncoder(w).Encode(m)
        }
    })

    err = http.ListenAndServe(*listenAddress, nil)
    if err != nil {
        log.Fatal(err)
    }
}

func createBatchPoint(m sensor.Measurement) (client.BatchPoints, error) {
    // Create a new point batch
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  *influxDbName,
        Precision: "s",
    })

    if err != nil {
        log.Println("failed to initialize batch points: ", err)
        return nil, err
    }

    // Set Tags
    var tags map[string]string = make(map[string]string)
    if *locationTag != "" {
        tags["location"] = *locationTag
    }
    if *roomTag != "" {
        tags["room"] = *roomTag
    }

    //Set Temperature
    tpt, err := client.NewPoint("temperature", tags, map[string]interface{}{"value": m.Temperature}, m.Time)
    if err != nil {
        log.Println("Error: ", err)
        return nil, err
    }
    bp.AddPoint(tpt)

    //Set Humidity
    hpt, err := client.NewPoint("humidity", tags, map[string]interface{}{"value": m.Humidity}, m.Time)
    if err != nil {
        log.Println("Error: ", err)
        return nil, err
    }
    bp.AddPoint(hpt)
    return bp, nil
}


func dbWriter(s *sensor.Sensor) {
    c, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: *influxDbUrl,
        Username: *influxDbUser,
        Password: *influxDbPassword,
    })

    if err != nil {
        log.Fatalln("failed to initialize db connection", err)
    }

    for {
        select {
            case m := <-s.Data:
                bp, err := createBatchPoint(m)
                if err != nil {
                    log.Println("could not create batch point", err)
                    continue
                }
                err = c.Write(bp)
                if err != nil {
                    log.Println("failed to write batch point to db: ", err)
                }
        }
    }
}

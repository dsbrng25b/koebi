package sensor

import (
    "sync"
    "time"
    "errors"
    "log"
    dht "github.com/d2r2/go-dht"
)

var MAX_MEASURE_TIME = 5

type Measurement struct {
    Temperature float32
    Humidity float32
    Time time.Time
}

type Sensor struct {
    sync.RWMutex
    pin int
    interval time.Duration
    currentMeasurement *Measurement
    sensorType dht.SensorType
    Data chan Measurement
}

func New(sensorType string, pin int, interval time.Duration) (*Sensor, error) {
    s := &Sensor{
        pin: pin, 
        interval: interval,
        Data: make(chan Measurement, 20),
    }

    switch sensorType {
        case "DHT11":
            s.sensorType = dht.DHT11
        case "DHT22":
            s.sensorType = dht.DHT22
        default:
            return nil, errors.New("invalied sensor type")
    }

    return s, nil
}

func (s *Sensor) GetLast() (m Measurement, err error){
    s.RLock()
    defer s.RUnlock()
    if s.currentMeasurement == nil {
        return Measurement{}, errors.New("no measurement available")
    }
    return *s.currentMeasurement, nil
}

func (s *Sensor) measure() (m *Measurement, err error){
    //measure code goes here
    t, h, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, s.pin, false, 5)

    if err != nil {
        return nil, err
    }

    m = &Measurement{
        Temperature: t, 
        Humidity: h,
        Time: time.Now(),
    }

    s.Lock()
    s.currentMeasurement = m
    s.Unlock()
    s.Data <-*m
    return m, nil
}

func (s *Sensor) Start() {
    for {
        _, err := s.measure()
        if err != nil {
            log.Println("measurement failed: ", err)
        }
        time.Sleep(s.interval)
    }
}


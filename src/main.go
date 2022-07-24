package main

import (
  "fmt"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "github.com/quhar/bme280"
  "golang.org/x/exp/io/i2c"
  "log"
  "net/http"
  "time"
)

var (
  sensorData = prometheus.NewGaugeVec(prometheus.GaugeOpts{
    Name: "sensor_data_from_raspberry_pi",
    Help: "Sensor data from raspberry pi",
  },
    []string{
      "key",
    },
  )
)

func init() {
  // Metrics have to be registered to be exposed:
  prometheus.MustRegister(sensorData)
}

// BME280 デバイスの初期化
func initBme280() (*bme280.BME280, error) {
  dev, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)

  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  bme280 := bme280.New(dev)
  err = bme280.Init()
  return bme280, err
}

// BME280 デバイスのセンサーからイベントを取得して表示
func outputSensorValues(bme *bme280.BME280) (err error) {
  temperature, pressure, humidity, err := bme.EnvData()
  if err != nil {
    fmt.Println(err)
    return
  }

  // 標準出力
  fmt.Printf("Temperature: %.2fC, Humidity: %.2f%%, Pressure: %.2fhpa \n", temperature, humidity, pressure)

  // Node exporter
  sensorData.With(prometheus.Labels{"key": "Temperature"}).Set(temperature)
  sensorData.With(prometheus.Labels{"key": "Humidity"}).Set(humidity)
  sensorData.With(prometheus.Labels{"key": "Pressure"}).Set(pressure)
  return
}

func setMetrics() {
  // インターバル
  const outputIntervalSec = 60

  // BME280 デバイスの初期化
  bme280, err := initBme280()
  if err != nil {
    fmt.Println(err)
    return
  }

  // 値取得間隔の設定
  ticker := time.NewTicker(outputIntervalSec * time.Second)

  // 終了したら Stop を実行する
  defer ticker.Stop()

  for {
    // チャネルの送受信操作の多重化
    select {
    // ticker を利用して、一定間隔で繰り返し実行
    case <-ticker.C:
      err = outputSensorValues(bme280)
      if err != nil {
        fmt.Println(err)
        return
      }
    }
  }
}

func main() {

  go setMetrics()

  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(":8080", nil))
}

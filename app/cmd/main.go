package main

import (
  "fmt"
  "github.com/quhar/bme280"
  "golang.org/x/exp/io/i2c"
  "os"
  "time"
)

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

  // ファイルへの書き込み
  file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    fmt.Println(err)
    return
  }
  defer file.Close()
  newLine := fmt.Sprintf("Temperature: %.2fC, Humidity: %.2f%%, Pressure: %.2fhpa", temperature, humidity, pressure)
  _, err = fmt.Fprintln(file, newLine)
  if err != nil {
    fmt.Println(err)
  }

  // 標準出力
  fmt.Printf("Temperature: %.2fC, Humidity: %.2f%%, Pressure: %.2fhpa \n", temperature, humidity, pressure)
  return
}

func main() {
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

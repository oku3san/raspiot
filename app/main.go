package main

import (
  "fmt"
  "github.com/quhar/bme280"
  "golang.org/x/exp/io/i2c"
)

func main() {
  dev, _ := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x76)
  bme := bme280.New(dev)
  bme.Init()
  temperature, pressure, humidity, _ := bme.EnvData()
  fmt.Println(temperature, pressure, humidity)
}

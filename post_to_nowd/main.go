package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/mschoch/go-bmp085"
)

var url = flag.String("url", "http://localhost:4793", "root url of nowd server")
var device = flag.String("device", "sensor", "name of device")
var elevation = flag.Float64("elevation", 0, "Elevation of sensor in meters")

func main() {

	flag.Parse()

	var i2cbus byte = 1
	d := new(bmp085.Device)
	err := d.Init(i2cbus)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// set mode
	d.SetMode(bmp085.ULTRAHIGHRES)

	temp, err := d.ReadTemp()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	pressure, err := d.ReadPressure()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	adjustedPressure := float64(pressure) / math.Pow(1-(*elevation/44330), 5.255)
	adjustedPressureHg := adjustedPressure / bmp085.TO_INCHES_MERCURY

	jsonData := map[string]interface{}{
		"t":  temp,
		"p":  adjustedPressureHg,
		"ts": time.Now().Format(time.RFC3339),
	}

	jsonOut, err := json.Marshal(jsonData)
	if err != nil {
		log.Fatalf("error marshaling JSON: %v", err)
	}
	jsonBuffer := bytes.NewBuffer(jsonOut)
	resp, err := http.Post(*url+"/"+*device, "application/json", jsonBuffer)
	if err != nil {
		log.Fatalf("error http post: %v", err)
	}
	log.Printf("resp: %s", resp.Status)
}

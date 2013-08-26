package main

import (
	"fmt"

	"github.com/mschoch/go-bmp085"
)

func main() {
	var i2cbus byte = 1
	d := new(bmp085.Device)
	err := d.Init(i2cbus)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	temp, err := d.ReadTemp()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Temp is %fC\n", temp)
}

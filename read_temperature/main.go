package main

import (
	"fmt"

	"github.com/mschoch/go-bmp085"
)

func main() {
	d := new(bmp085.Device)
	err := d.Init(1)
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

//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/mschoch/go-bmp085"
)

var elevation = flag.Float64("elevation", 0, "Elevation of sensor in meters")

func main() {
	flag.Parse()

	var i2cbus byte = 1
	d := new(bmp085.Device)
	err := d.Init(i2cbus)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	pressure, err := d.ReadPressure()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("Pressure is %dpa\n", pressure)

	if *elevation != 0.0 {
		adjustedPressure := float64(pressure) / math.Pow(1-(*elevation/44330), 5.255)
		fmt.Printf("Sea-level adjusted pressure is %fpa", adjustedPressure)
	}
}

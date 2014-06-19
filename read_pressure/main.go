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

var raw = flag.Bool("raw", false, "print raw values not adjusted to sea-level")
var sealevel = flag.Bool("sealevel", true, "print sea-level adjusted values")
var hg = flag.Bool("hg", false, "print value in inches mercury")
var pa = flag.Bool("pa", true, "print value in pascals")
var mode = flag.Int("mode", int(bmp085.STANDARD), "mode: 0-ultralowpower, 1-standard, 2-highres, 3-ultrahighres")
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

	// set mode
	d.SetMode(byte(*mode))

	pressure, err := d.ReadPressure()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	if *raw {
		fmt.Printf("Raw pressure:\n")
		if *pa {
			fmt.Printf("\t%d pa\n", pressure)
		}
		if *hg {
			pressureHg := float64(pressure) / bmp085.TO_INCHES_MERCURY
			fmt.Printf("\t%f in\n", pressureHg)
		}
	}

	if *sealevel {
		fmt.Printf("Sea-level adjusted pressure:\n")
		adjustedPressure := float64(pressure) / math.Pow(1-(*elevation/44330), 5.255)
		if *pa {
			fmt.Printf("\t%f pa\n", adjustedPressure)
		}

		if *hg {
			adjustedPressureHg := adjustedPressure / bmp085.TO_INCHES_MERCURY
			fmt.Printf("\t%f in\n", adjustedPressureHg)
		}
	}

}

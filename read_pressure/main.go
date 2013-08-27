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

	d.ReadPressure()
}

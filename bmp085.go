//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package bmp085

import (
	"bytes"
	"encoding/binary"
	"time"

	"bitbucket.org/gmcbay/i2c"
)

const I2C_ADDR = 0x77

const CAL_AC1 byte = 0xAA
const CAL_AC2 byte = 0xAC
const CAL_AC3 byte = 0xAE
const CAL_AC4 byte = 0xB0
const CAL_AC5 byte = 0xB2
const CAL_AC6 byte = 0xB4
const CAL_B1 byte = 0xB6
const CAL_B2 byte = 0xB8
const CAL_MB byte = 0xBA
const CAL_MC byte = 0xBC
const CAL_MD byte = 0xBE

const CONTROL = 0xF4
const TEMP_DATA = 0xF6
const PRESSURE_DATA = 0xF6
const READ_TEMP_CMD = 0x2E
const READ_PRESSURE_CMD = 0x34

type Device struct {
	bus    *i2c.I2CBus
	busNum byte
	addr   byte
	ac1    int16
	ac2    int16
	ac3    int16
	ac4    uint16
	ac5    uint16
	ac6    uint16
	b1     int16
	b2     int16
	mb     int16
	mc     int16
	md     int16
}

func (d *Device) Init(busNum byte) (err error) {
	return d.InitCustomAddr(I2C_ADDR, busNum)
}

func (d *Device) InitCustomAddr(addr, busNum byte) (err error) {
	if d.bus, err = i2c.Bus(busNum); err != nil {
		return
	}

	d.busNum = busNum
	d.addr = addr

	err = d.readCalibration()

	return
}

func (d *Device) readCalibration() (err error) {
	var data []byte

	calibrations := []struct {
		addr byte
		dest interface{}
	}{
		{CAL_AC1, &d.ac1},
		{CAL_AC2, &d.ac2},
		{CAL_AC3, &d.ac3},
		{CAL_AC4, &d.ac4},
		{CAL_AC5, &d.ac5},
		{CAL_AC6, &d.ac6},
		{CAL_B1, &d.b1},
		{CAL_B2, &d.b2},
		{CAL_MB, &d.mb},
		{CAL_MC, &d.mc},
		{CAL_MD, &d.md},
	}

	for _, calibration := range calibrations {
		if data, err = d.bus.ReadByteBlock(d.addr, calibration.addr, 2); err != nil {
			return
		}
		p := bytes.NewBuffer(data)
		err = binary.Read(p, binary.BigEndian, calibration.dest)
		if err != nil {
			return
		}
	}

	return
}

func (d *Device) readUncalibratedTemp() (temp int16, err error) {
	var data []byte
	if err = d.bus.WriteByte(d.addr, CONTROL, READ_TEMP_CMD); err != nil {
		return
	}
	time.Sleep(5 * time.Millisecond)
	if data, err = d.bus.ReadByteBlock(d.addr, TEMP_DATA, 2); err != nil {
		return
	}
	p := bytes.NewBuffer(data)
	err = binary.Read(p, binary.BigEndian, &temp)
	return
}

func (d *Device) ReadTemp() (temp float32, err error) {
	var utraw int16
	if utraw, err = d.readUncalibratedTemp(); err != nil {
		return
	}
	ut := int32(utraw)
	ac6 := int32(d.ac6)
	ac5 := int32(d.ac5)
	x1 := ((ut - ac6) * ac5) >> 15

	mc := int32(d.mc)
	md := int32(d.md)
	x2 := (mc << 11) / (x1 + md)
	b5 := x1 + x2
	t := (b5 + 8) / 16
	temp = float32(t) / 10.0
	return
}

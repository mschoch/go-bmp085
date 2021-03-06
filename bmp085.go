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

const TO_INCHES_MERCURY = 3386.389

const I2C_ADDR = 0x77

const (
	ULTRALOWPOWER byte = iota
	STANDARD
	HIGHRES
	ULTRAHIGHRES
)

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
	mode   byte
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

func (d *Device) readUncalibratedPressure() (pressure int32, err error) {
	var data []byte
	if err = d.bus.WriteByte(d.addr, CONTROL, READ_PRESSURE_CMD+(d.mode<<6)); err != nil {
		return
	}
	switch d.mode {
	case ULTRALOWPOWER:
		time.Sleep(5 * time.Millisecond)
	case STANDARD:
		time.Sleep(8 * time.Millisecond)
	case HIGHRES:
		time.Sleep(14 * time.Millisecond)
	case ULTRAHIGHRES:
		time.Sleep(26 * time.Millisecond)
	}

	if data, err = d.bus.ReadByteBlock(d.addr, PRESSURE_DATA, 3); err != nil {
		return
	}
	zpadData := make([]byte, 4)
	copy(zpadData[1:], data)
	p := bytes.NewBuffer(zpadData)
	err = binary.Read(p, binary.BigEndian, &pressure)
	pressure = pressure >> (8 - d.mode)

	return
}

func (d *Device) ReadPressure() (p int32, err error) {

	var utraw int16
	if utraw, err = d.readUncalibratedTemp(); err != nil {
		return
	}

	var upraw int32
	if upraw, err = d.readUncalibratedPressure(); err != nil {
		return
	}

	// create larger vars to avoid overflowing
	ut := int32(utraw)
	ac1 := int32(d.ac1)
	ac2 := int32(d.ac2)
	ac3 := int32(d.ac3)
	ac4 := uint32(d.ac4)
	ac6 := int32(d.ac6)
	ac5 := int32(d.ac5)
	b1 := int32(d.b1)
	b2 := int32(d.b2)
	mc := int32(d.mc)
	md := int32(d.md)

	//calculate temp
	x1 := ((ut - ac6) * ac5) >> 15
	x2 := (mc << 11) / (x1 + md)
	b5 := x1 + x2
	//t := (b5 + 8) / 16

	// calculate pressure
	b6 := b5 - 4000
	x1 = (b2 * ((b6 * b6) >> 12)) >> 11
	x2 = (ac2 * b6) >> 11
	x3 := x1 + x2
	b3 := (((ac1*4 + x3) << d.mode) + 2) / 4

	x1 = (ac3 * b6) >> 13
	x2 = (b1 * ((b6 * b6) >> 12)) >> 16
	x3 = ((x1 + x2) + 2) >> 2
	var tmpa = uint32(x3 + 32768)
	b4 := (ac4 * tmpa) >> 15
	var tmpb = uint32(upraw - b3)
	var tmpc = uint32(50000 >> d.mode)
	b7 := tmpb * tmpc
	p = int32((b7 / b4) / 2)
	if b7 < 0x80000000 {
		p = int32((b7 * 2) / b4)
	}
	x1 = (p >> 8) * (p >> 8)
	x1 = (x1 * 3038) >> 16
	x2 = (-7367 * p) >> 16
	p = p + ((x1 + x2 + int32(3791)) >> 4)
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

func (d *Device) SetMode(mode byte) {
	d.mode = mode
}

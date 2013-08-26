## go-bmp085

A go library to interact with the [BMP085](http://dlnmh9ip6v2uc.cloudfront.net/datasheets/Sensors/Pressure/BST-BMP085-DS000-06.pdf) Temperature/Pressure Sensor.  

NOTE:  This has only been tested on a Raspberry Pi.

### Usage

    import ("github.com/mschoch/go-bmp085")

    // choose which i2c bus to use
    var i2cbus byte = 1

    // create the device
    d := new(bmp085.Device)

    // initialize it
    _ := d.Init(i2cbus)

    // read the temperature
    temp, _ := d.ReadTemp()

### Examples

There is one example app showing how to read the temperature. (you may need to be root to read the i2c bus)

    $ sudo ./read_temperature 
    Temp is 20.900000C


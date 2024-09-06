# Client Devices

## Raspberry Pi I2C

* Don't forget to enable I2C using raspi-config
* This project only uses the I2C "1" Bus

### Debugging I2C

* Install `i2c-tools` and use `i2cdetect -y 1` to verify the address of connected I2C devices

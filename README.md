# pi-vitrine
Display your pi

## Installing Go on a Raspberry Pi

Follow the normal instructions for downloading and installing Gol ([link](https://go.dev/doc/install)). 
You just have to make sure you download the correct file.

- If the OS on the Raspberry Pi is 64 bit, you should download the ...linux-arm64.tar.gz file
- If the OS on the Raspberry Pi is 32 bit, you should download the ...linux-armv6l.tar.gz file

You can use `uname -a` to figure out the OS info on the Raspberry Pi if you ever forget.

## Host Server

The host server provides endpoints for storing data sent by client devices.

See [here](db/README.md) for postgres database setup.

## Client Devices

Client devices read sensor data and send it to the host server.

See [here](cmd/client/README.md) for client setup instructions.

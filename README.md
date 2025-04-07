# pi-vitrine
An IoT client-server system, designed for running at home on Raspberry Pi.

## Installing Go on a Raspberry Pi

Follow the normal instructions for downloading and installing Go found [here](https://go.dev/doc/install). 

You just have to make sure you download the correct file.

- If the OS on the Raspberry Pi is 64 bit, you should download the ...linux-arm64.tar.gz file
- If the OS on the Raspberry Pi is 32 bit, you should download the ...linux-armv6l.tar.gz file

You can use `uname -a` to figure out the OS info on the Raspberry Pi if you ever forget.

## Host Server

The host server provides a web UI for the user to register client devices and visualize data.
It also provides HTTP endpoints for storing data sent by client devices.

See [here](db/README.md) for postgres database setup.

Build and run the host server from the [cmd/server/](cmd/server/) package.

* There is a `--help` option to explain the different CLI parameters.


## Client Raspberry Pi Device

A client Raspberry Pi device reads sensor data and sends it to the host server over HTTP. It also provides a web UI for displaying client sensor data.

See [here](cmd/client/README.md) for client setup help.

Build and run the client from the [cmd/client/](cmd/client/) package.

* There is a `--help` option to explain the different CLI parameters.

# pi-vitrine
An IoT client-server system, designed for running at home on Raspberry Pi.

## Installing Go on a Raspberry Pi

Follow the normal instructions for downloading and installing Go found [here](https://go.dev/doc/install). 

You just have to make sure you download the correct file.

- If the OS on the Raspberry Pi is 64 bit, you should download the ...linux-arm64.tar.gz file
- If the OS on the Raspberry Pi is 32 bit, you should download the ...linux-armv6l.tar.gz file

You can use `uname -a` to figure out the OS info on the Raspberry Pi if you ever forget.

## Host Server

The host server provides a web UI as well as HTTP endpoints for storing data sent by client devices.

See [here](db/README.md) for postgres database setup.

Run the host server from the [cmd/server/](cmd/server/) package.


## Client Device

A client device reads sensor data and sends it to the host server over HTTP. It also provides a web UI for configuring the client.

See [here](cmd/client/README.md) for client setup help.

Run the client from the [cmd/client/](cmd/client/) package.

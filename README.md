# github.com/SkYNewZ/yeelight-control-over-process

> Project is work is progress.

`yeelight-control-over-process` is a background program that controls your Yeelights as you have configured it by listening running processes on your computer.

It works on all distribution supported by https://github.com/mitchellh/go-ps

## Getting started

1. Install this binary
2. Create a configuration file matching [`config.yaml`](/config.yaml). All settings are documented into this file
3. Run the program

When a process on your machine is detected as present or absent, the actions defined in the configuration will take place.

For example:

- if I start `VLC`, set my desktop light to red
- if I close `VLC`, turn off this light

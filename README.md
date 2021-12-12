# periph - Peripherals I/O in Go

Documentation is at https://periph.io

Join us for a chat on
[gophers.slack.com/messages/periph](https://gophers.slack.com/messages/periph),
get an [invite here](https://invite.slack.golangbridge.org/).

[![mascot](https://raw.githubusercontent.com/periph/website/master/site/static/img/periph-mascot-280.png)](https://periph.io/)

[![PkgGoDev](https://pkg.go.dev/badge/periph.io/x/cmd)](https://pkg.go.dev/periph.io/x/cmd)


# cmd - ready-to-use executables

This directory contains directly usable tools installable via:

```
go get periph.io/x/cmd/...
```


## Push

If you prefer to build on your workstation and push the binaries to the micro
computer, install `push` from [periph.io/x/bootstrap](
https://github.com/periph/bootstrap) to cross compile and efficiently push via
rsync:

```
go get -u periph.io/x/bootstrap/cmd/push
push -host pi@raspberrypi periph.io/x/cmd/...
```


## Recommended first use

Try first `periph-info`. It will print out if any driver failed to run, for
example if you have to run as root to access certain drivers.

Then run `headers-list` to list all the headers on your board and confirm that
you get the expected output. If your board is missing, you can [contribute
it](https://periph.io/project/contributing/).


## Devices

- [apa102](apa102): Writes to a LED strip of APA-102 (sometimes called Dotstar).
  Can show an image animating on the Y axis.
- [bmxx80](bmxx80): Reads the temperature, pressure and humidity off a
  bmp180/bme280/bmp280. Humidity sensing is only supported on bme280.
- [cap1xxx](cap1xxx): Reads the capacitive sensor family.
- [ir](ir): Reads codes (button presses) on an InfraRed remote sensor.
- [led](led): Reads the state of on-board LEDs.
- [ssd1306](ssd1306): Writes text, an image or an animated GIF to an OLED
  display.
- [tm1637](tm1637): Writes to a segment digits display.


## Buses

- [gpio-list](gpio-list): Looking for the GPIO pins per functionality?
  Prints the state of each GPIO pin.
- [gpio-read](gpio-read): Read the input value of a GPIO pin and change
  input resistor.
- [gpio-write](gpio-write): Change the output value of a GPIO pin.
- [headers-list](headers-list): Pinrts the location of the pin on the header to
  connect your GPIO. This is the perfect tool to know where to connect the
  wires.
- [i2c-io](i2c-io): Reads and/or writes to an I²C device.
- [i2c-list](i2c-list): Lists which I²C buses are enabled and where the pins
  are.
- [spi-io](spi-io): Reads and/or writes to an SPI device.
- [spi-list](spi-list): Lists which SPI ports are enabled and where the pins
  are.


## Other

- [periph-info](periph-info): Lists which periph drivers loaded and which
  failed.
- [periph-smoketest](periph-smoketest): Runs one of the smoke test for the
  drivers. The smoke test differs from unit tests as they require real hardware
  to confirm that the driver being tested works.
- [videosink](videosink): Demonstrates how to provide display contents to HTTP
  clients.


## Troubleshooting

Having trouble getting the tools to run? File [an
issue](https://github.com/periph/cmd/issues) or better visit the [Slack
channel](https://gophers.slack.com/messages/periph/). You can get an [invite
here](https://invite.slack.golangbridge.org/).


## Authors

`periph` was initiated with ❤️️ and passion by [Marc-Antoine
Ruel](https://github.com/maruel). The full list of contributors is in
[AUTHORS](https://github.com/periph/cmd/blob/main/AUTHORS) and
[CONTRIBUTORS](https://github.com/periph/cmd/blob/main/CONTRIBUTORS).


## Disclaimer

This is not an official Google product (experimental or otherwise), it
is just code that happens to be owned by Google.

This project is not affiliated with the Go project.

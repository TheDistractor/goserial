/*
Goserial is a simple go package to allow you to read and write from
the serial port as a stream of bytes.

It aims to have the same API on all platforms, including windows.  As
an added bonus, the windows package does not use cgo, so you can cross
compile for windows from another platform.  Unfortunately goinstall
does not currently let you cross compile so you will have to do it
manually:

 GOOS=windows make clean install

Currently there is very little in the way of configurability.  You can
set the baud rate.  Then you can Read(), Write(), or Close() the
connection.  Read() will block until at least one byte is returned.
Write is the same.  There is currently no exposed way to set the
timeouts, though patches are welcome.

Currently ports are opened with 8 data bits, 1 stop bit, no parity, no hardware
flow control, and no software flow control by default.  This works fine for
many real devices and many faux serial devices including usb-to-serial
converters and bluetooth serial ports.

You may Read() and Write() simulantiously on the same connection (from
different goroutines).

Example usage:

  package main

  import (
        "github.com/tarm/goserial"
        "log"
  )

  func main() {
        c := &serial.Config{Name: "COM5", Baud: 115200}
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }

        n, err := s.Write([]byte("test"))
        if err != nil {
                log.Fatal(err)
        }

        buf := make([]byte, 128)
        n, err = s.Read(buf)
        if err != nil {
                log.Fatal(err)
        }
        log.Print("%q", buf[:n])
  }
*/
package goserial

import (
	"errors"
	"io"
)

var (
	ErrConfigStopBits = errors.New("goserial config: bad number of stop bits")
	ErrConfigByteSize = errors.New("goserial config: bad byte size")
	ErrConfigParity   = errors.New("goserial config: bad parity")
)

type ParityMode byte

const (
	ParityNone = ParityMode(iota)
	ParityEven
	ParityOdd
)

type ByteSize byte

const (
	Byte8 = ByteSize(iota)
	Byte5
	Byte6
	Byte7
)

type StopBits byte

const (
	StopBits1 = StopBits(iota)
	StopBits2
)

// Config contains the information needed to open a serial port.
//
// Currently few options are implemented, but more may be added in the
// future (patches welcome), so it is recommended that you create a
// new config addressing the fields by name rather than by order.
//
// For example:
//
//    c0 := &serial.Config{Name: "COM45", Baud: 115200}
// or
//    c1 := new(serial.Config)
//    c1.Name = "/dev/tty.usbserial"
//    c1.Baud = 115200
//
type Config struct {
	Name string
	Baud int

	Size     ByteSize
	Parity   ParityMode
	StopBits StopBits

	// RTSFlowControl bool
	// DTRFlowControl bool
	// XONFlowControl bool

	CRLFTranslate bool // Ignored on Windows.
	// TimeoutStuff int
	ReadTimeout uint32
}

func (c *Config) check() error {
	switch c.Size {
	case Byte5, Byte6, Byte7, Byte8:
	default:
		return ErrConfigByteSize
	}

	switch c.StopBits {
	case StopBits1, StopBits2:
	default:
		return ErrConfigParity
	}

	switch c.Parity {
	case ParityNone, ParityEven, ParityOdd:
	default:
		return ErrConfigParity
	}

	return nil
}

// OpenPort opens a serial port with the specified configuration
func OpenPort(c *Config) (io.ReadWriteCloser, error) {
	if err := c.check(); err != nil {
		return nil, err
	}

	return openPort(c.Name, c)
}

// func Flush()

// func SendBreak()

// func RegisterBreakHandler(func())

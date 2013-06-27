package main

import (
	"fmt"
	"os"
)

const (
	RED_PIN   = "18"
	GREEN_PIN = "23"
	BLUE_PIN  = "24"
)

func init() {
	exportPin(RED_PIN)
	exportPin(GREEN_PIN)
	exportPin(BLUE_PIN)

	setPinDirection(RED_PIN)
	setPinDirection(GREEN_PIN)
	setPinDirection(BLUE_PIN)

}

type LED struct {
	redPin *os.File
	redOn  bool

	greenPin *os.File
	greenOn  bool

	bluePin *os.File
	blueOn  bool
}

func NewLED() *LED {
	led := &LED{}
	var err error

	led.redPin, err = os.OpenFile("/sys/class/gpio/gpio"+RED_PIN+"/value", os.O_WRONLY|os.O_TRUNC, 755)
	if err != nil {
		panic("Failed to open GPIO pin")
	}

	led.greenPin, err = os.OpenFile("/sys/class/gpio/gpio"+GREEN_PIN+"/value", os.O_WRONLY|os.O_TRUNC, 755)
	if err != nil {
		panic("Failed to open GPIO pin")
	}

	led.bluePin, err = os.OpenFile("/sys/class/gpio/gpio"+BLUE_PIN+"/value", os.O_WRONLY|os.O_TRUNC, 755)
	if err != nil {
		panic("Failed to open GPIO pin")
	}

	return led
}

func (l *LED) Red(on bool) {
	if on && l.redOn {
		return
	} else if on {
		l.redPin.Write([]byte("1"))
		l.redOn = true
	} else {
		l.redPin.Write([]byte("0"))
		l.redOn = false
	}

	return
}

func (l *LED) Green(on bool) {
	if on && l.greenOn {
		return
	} else if on {
		l.greenPin.Write([]byte("1"))
		l.greenOn = true
	} else {
		l.greenPin.Write([]byte("0"))
		l.greenOn = false
	}

	return
}

func (l *LED) Blue(on bool) {
	if on && l.blueOn {
		return
	} else if on {
		l.bluePin.Write([]byte("1"))
		l.blueOn = true
	} else {
		l.bluePin.Write([]byte("0"))
		l.blueOn = false
	}

	return
}

func (l *LED) Off() {
	l.Red(false)
	l.Green(false)
	l.Blue(false)
}

func (l *LED) Shutdown() {
	fmt.Println("Shutting Down")
	l.Off()

	exportPin("17")
	exportPin("21")
	exportPin("22")

	l.redPin.Close()
	l.greenPin.Close()
	l.bluePin.Close()
}

func setPinDirection(pin string) {
	dir, err := os.OpenFile("/sys/class/gpio/gpio"+pin+"/direction", os.O_WRONLY|os.O_TRUNC, 755)

	if err != nil {
		panic("Failed to set GPIO pin direction: " + pin)
	}

	defer dir.Close()

	dir.Write([]byte("out"))
}

func exportPin(pin string) {
	exp, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY|os.O_TRUNC, 755)

	if err != nil {
		panic("Failed to export GPIO pin: " + pin)
	}

	defer exp.Close()
	exp.Write([]byte(pin))
}

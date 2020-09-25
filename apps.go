package main

import (
	"github.com/SkYNewZ/go-yeelight"
	log "github.com/sirupsen/logrus"
)

var err error
var strip1 *yeelight.Yeelight
var strip2 *yeelight.Yeelight
var todos []*processConfig

func init() {
	strip1, err = yeelight.New("192.168.1.44", "Strip bureau")
	checkError(err)
	log.Printf("[%s] Connected\n", strip1.Name)

	strip2, err = yeelight.New("192.168.1.15", "Strip salon")
	checkError(err)
	log.Printf("[%s] Connected\n", strip2.Name)

	// Define what to do with this light
	todos = []*processConfig{
		{
			processNames: []string{"notepad.exe", "calculator.exe"},
			light:        strip1,
			onFound:      makeItGreen,
			onNotFound:   makeItRed,
		},
		{
			processNames: []string{"notepad.exe", "calculator.exe"},
			light:        strip2,
			onFound:      makeItGreen,
			onNotFound:   makeItRed,
		},
	}
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/SkYNewZ/go-yeelight"
	"github.com/mitchellh/go-ps"
	log "github.com/sirupsen/logrus"
)

type command int

const (
	on command = iota
	off
	rgb
)

type action struct {
	done   bool
	action command
	args   interface{}
}

type processConfig struct {
	processNames []string
	light        *yeelight.Yeelight
	onFound      action
	onNotFound   action
	mux          sync.Mutex
}

func (t *processConfig) toogle(found bool) {
	t.onFound.done = found
	t.onNotFound.done = !found
}

var (
	red       [3]int = [3]int{255, 0, 0}
	green     [3]int = [3]int{0, 255, 0}
	blue      [3]int = [3]int{0, 0, 255}
	processes []ps.Process
	errs      []*error

	makeItRed action = action{
		action: rgb,
		args:   red,
	}
	makeItGreen action = action{
		action: rgb,
		args:   green,
	}
	makeItBlue action = action{
		action: rgb,
		args:   blue,
	}
)

// Generic function to make command on lights
func genericFunc(a command, b *yeelight.Yeelight, args interface{}) error {
	switch a {
	case rgb:
		if colors, ok := args.([3]int); ok {
			b.RGB(colors[0], colors[1], colors[2])
		} else {
			return fmt.Errorf("Invalid colors")
		}
	case on:
		return b.TurnOn()
	case off:
		return b.TurnOff()
	}

	return nil
}

// Search if wanted process is currently running
func searchingMatchingProcess(t *processConfig) (bool, string) {
	for _, searchedProcess := range t.processNames {
		for _, p := range processes {
			found, name := strings.Contains(strings.ToLower(p.Executable()), strings.ToLower(searchedProcess)), searchedProcess
			if found {
				return found, name
			}
		}
	}

	return false, ""
}

// Generic function to throw error
func checkError(err error) {
	if err == nil {
		return
	}

	errs = append(errs, &err)
	log.Errorf("%s", err)

	if len(errs) > 5 {
		log.Errorln("Too many errors, exit")
		os.Exit(1)
	}
}

// Listing current process
func startProcessesProcess(ticker *time.Ticker) {
	for ; true; <-ticker.C {
		p, err := ps.Processes()
		checkError(err)
		processes = p
	}
}

func setupExitHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\rExiting")
		processes = nil
		errs = nil
		os.Exit(0)
	}()
}

func main() {
	// Get control+c event
	setupExitHandler()

	// Start process scanning
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go startProcessesProcess(ticker)
	time.Sleep(time.Second * 2)

	// Main program
	for {
		// For each wanted program to be the trigger
		for _, v := range todos {

			// Search if program is running
			found, name := searchingMatchingProcess(v)

			var action command
			var args interface{}

			// If program found, do action and lock the device to avoid loop on it
			if found && !v.onFound.done {
				log.Printf("[%s] Processes '%s' found, applying changes\n", v.light.Name, name)
				v.toogle(found)
				action = v.onFound.action
				args = v.onFound.args
			}

			// If process not found, do action and unlock device
			if !found && !v.onNotFound.done {
				log.Printf("[%s] Any of '%s' processes not found, applying changes\n", v.light.Name, strings.Join(v.processNames, ", "))
				v.toogle(found)
				action = v.onNotFound.action
				args = v.onNotFound.args
			}

			// If not power on, but wanted to manage light
			if action == rgb {
				powerOn, err := v.light.IsPowerOn()
				checkError(err)
				if !powerOn {
					log.Printf("[%s] Powering on\n", v.light.Name)
					checkError(v.light.TurnOn())
				}
			}

			checkError(genericFunc(action, v.light, args))
		}
	}
}

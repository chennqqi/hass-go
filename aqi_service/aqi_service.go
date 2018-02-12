package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/jurgen-kluft/hass-go/aqi"
	"github.com/jurgen-kluft/hass-go/state"
)

func main() {
	states := state.New()
	aqiInstance, _ := aqi.New()

	states.SetTimeState("time", "now", time.Now())
	aqiInstance.Process(states)

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for t := range ticker.C {
			states.SetTimeState("time", "now", t)
			aqiInstance.Process(states)
		}
	}()

	for true {
		fmt.Print("\n(P)rint or (E)xit? ")
		reader := bufio.NewReader(os.Stdin)
		c, _ := reader.ReadByte()
		if c == 'P' {
			states.Print()
		}
		if c == 'E' {
			break
		}
	}
	ticker.Stop()
}

package hass

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/jurgen-kluft/hass-go/state"
)

func UpdateHttpSensor(name string, state string, unit string, friendlyname string) {
	url := "http://10.0.0.22:8123/api/states/sensor." + name
	json := "'{\"state\": \"$(STATE)\", \"attributes\": {\"unit_of_measurement\": \"$(UNIT)\", \"friendly_name\": \"$(FRIENDLY_NAME)\"}}'"
	json = strings.Replace(json, "$(STATE)", state, 1)
	json = strings.Replace(json, "$(UNIT)", unit, 1)
	json = strings.Replace(json, "$(FRIENDLY_NAME)", friendlyname, 1)

	resp, err := http.Post(url, "application/json", bytes.NewBufferString(json))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

type Instance struct {
}

func (c *Instance) Process(states *state.Domain) {
	sensors := states.Get("hass")
	for sn, sv := range sensors.Strings {
		fname := states.GetStringState("sensor", sn+".descr", "?")
		unit := states.GetStringState("sensor", sn+".unit", "")
		UpdateHttpSensor(sn, sv, unit, fname)
	}
}

package aqi

import "encoding/json"

func unmarshalcaqi(data []byte) (*Caqi, error) {
	r := &Caqi{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Caqi) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Caqi struct {
	URL string `json:"url"`
}

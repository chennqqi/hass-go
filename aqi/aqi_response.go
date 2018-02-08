// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCaqiResponse(bytes)
//    bytes, err = r.Marshal()

package aqi

import "encoding/json"

func unmarshalCaqiResponse(data []byte) (CaqiResponse, error) {
	var r CaqiResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CaqiResponse) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CaqiResponse struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	Aqi          int64         `json:"aqi"`
	Idx          int64         `json:"idx"`
	Attributions []Attribution `json:"attributions"`
	City         City          `json:"city"`
	Dominentpol  string        `json:"dominentpol"`
	Iaqi         Iaqi          `json:"iaqi"`
	Time         Time          `json:"time"`
}

type Attribution struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type City struct {
	Name string    `json:"name"`
	URL  string    `json:"url"`
	Geo  []float64 `json:"geo"`
}

type Iaqi struct {
	Co   Co `json:"co"`
	D    D  `json:"d"`
	H    D  `json:"h"`
	No2  Co `json:"no2"`
	O3   Co `json:"o3"`
	P    D  `json:"p"`
	Pm10 D  `json:"pm10"`
	Pm25 D  `json:"pm25"`
	So2  Co `json:"so2"`
	T    D  `json:"t"`
	W    D  `json:"w"`
	Wd   D  `json:"wd"`
}

type Co struct {
	V float64 `json:"v"`
}

type D struct {
	V int64 `json:"v"`
}

type Time struct {
	S  string `json:"s"`
	Tz string `json:"tz"`
	V  int64  `json:"v"`
}

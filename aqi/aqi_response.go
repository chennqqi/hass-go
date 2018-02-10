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
	Co   float64 `json:"co"`
	D    float64 `json:"d"`
	H    float64 `json:"h"`
	No2  float64 `json:"no2"`
	O3   float64 `json:"o3"`
	P    float64 `json:"p"`
	Pm10 float64 `json:"pm10"`
	Pm25 float64 `json:"pm25"`
	So2  float64 `json:"so2"`
	T    float64 `json:"t"`
	W    float64 `json:"w"`
	Wd   float64 `json:"wd"`
}

type Time struct {
	S  string `json:"s"`
	Tz string `json:"tz"`
	V  int64  `json:"v"`
}

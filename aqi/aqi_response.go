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
	Idx          int64         `json:"idx"`
	Aqi          int64         `json:"aqi"`
	Time         Time          `json:"time"`
	City         City          `json:"city"`
	Attributions []Attribution `json:"attributions"`
	Iaqi         Iaqi          `json:"iaqi"`
}

type Attribution struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type City struct {
	Name string   `json:"name"`
	URL  string   `json:"url"`
	Geo  []string `json:"geo"`
}

type Iaqi struct {
	Pm25 H `json:"pm25"`
	O3   H `json:"o3"`
	No2  H `json:"no2"`
	T    H `json:"t"`
	P    H `json:"p"`
	H    H `json:"h"`
}

type H struct {
	V int64 `json:"v"`
}

type Time struct {
	V  int64  `json:"v"`
	S  string `json:"s"`
	Tz string `json:"tz"`
}

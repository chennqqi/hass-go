// To parse and unparse this JSON data, add this code to your project and do:
//
//    r, err := UnmarshalCcalendar(bytes)
//    bytes, err = r.Marshal()

package calendar

import "encoding/json"

func unmarshalccalendar(data []byte) (*ccalendar, error) {
	r := &ccalendar{}
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ccalendar) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ccalendar struct {
	calendars []ccal   `json:"calendars"`
	event     []cevent `json:"event"`
}

type ccal struct {
	name string `json:"name"`
	url  string `json:"url"`
}

type cevent struct {
	calendar string   `json:"calendar"`
	domain   string   `json:"domain"`
	name     string   `json:"name"`
	state    string   `json:"state"`
	typeof   string   `json:"typeof"`
	values   []string `json:"values"`
}

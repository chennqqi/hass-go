package timeofday

import (
	"encoding/json"
	"fmt"
	"time"
)

func unmarshalctimeofday(data []byte) (*Ctimeofday, error) {
	r := &Ctimeofday{}
	r.Weekday = map[string][]Ctod{}
	err := json.Unmarshal(data, r)
	return r, err
}

func (r *Ctimeofday) marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Ctimeofday struct {
	Weekday map[string][]Ctod `json:"weekday"`
}

func (t *Ctimeofday) print() {
	for k, v := range t.Weekday {
		fmt.Printf("%s\n", k)
		for _, tod := range v {
			fmt.Printf("%s from %v \n", tod.Name, tod.From)
		}
	}
}

type Ctime struct {
	time.Time
}

func (ct *Ctime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse("15:04:05", s)
	ct.Time = t
	return err
}

type Ctod struct {
	Name string `json:"name"`
	From Ctime  `json:"from"`
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type stati struct {
	f             string
	Stadt         status
	OtzThueringen status
	OtzWeltweit   status
}

type status struct {
	Timestamp int64
	Count     int
}

func load(f string) (si stati) {
	si.f = f
	j, err := ioutil.ReadFile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(j, &si)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func (si stati) save() {
	j, err := json.Marshal(si)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(si.f, j, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

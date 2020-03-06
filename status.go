package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
)

type status struct {
	filename  string
	timestamp int64
	count     int
}

func (s *status) load(filename string) {
	s.filename = filename
	content, err := ioutil.ReadFile(s.filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := bytes.Split(content, []byte("\n"))
	s.timestamp, err = strconv.ParseInt(string(lines[0]), 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	count, err := strconv.ParseInt(string(lines[1]), 10, 0)
	if err != nil {
		fmt.Println(err)
	}
	s.count = int(count)
}

func (s *status) save() {
	err := ioutil.WriteFile(s.filename, []byte(fmt.Sprintf("%d\n%d", s.timestamp, s.count)), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

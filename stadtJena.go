package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/goodsign/monday"
)

const defaultStatus = "Für Jena liegen aktuell keine Meldungen über Erkrankungen vor. "

func stadtJena(filename string) {
	s.load(filename)
	c := colly.NewCollector()
	c.OnHTML(".content-inner--main > div > div:first-child > .paragraph--type--text > .text-formatted:first-child", stadtNeuigkeiten)
	c.Visit("http://jena.de/corona")
	s.save()
}

func stadtNeuigkeiten(h *colly.HTMLElement) {
	var statusElement, standElement *goquery.Selection
	h.DOM.Find("p").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Für Jena ") {
			statusElement = s
		}
		if strings.Contains(s.Text(), "Stand: ") {
			standElement = s
		}
	})
	if (statusElement == nil) || (standElement == nil) {
		html, _ := h.DOM.Html()
		fmt.Printf("Kann '%s' nicht auslesen\n", html)
		return
	}
	statusInfo := getStatus(statusElement)
	count := getCount(statusElement)
	timestamp := getTimestamp(standElement)
	zeitpunkt := monday.Format(timestamp, "2. January 15:04 Uhr", "de_DE")

	if s.count == 0 && statusInfo != defaultStatus {
		sendSignal("*WARNUNG - CORONA FALL* (%s)\n%s", zeitpunkt, statusInfo)
		s.count = 1
	} else if count > s.count {
		sendSignal("*Statusänderung* (%s)\nGestiegen von %d auf %d\n%s", zeitpunkt, s.count, count, statusInfo)
		s.count = count
	} else if timestamp.Unix() > s.timestamp {
		sendSignal("Aktualisierung vom %s\n%s", zeitpunkt, statusInfo)
		s.timestamp = timestamp.Unix()
	}
}

func getStatus(s *goquery.Selection) string {
	html, err := s.Children().First().Html()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return textize(html)
}

func getCount(s *goquery.Selection) int {
	number := s.Find("strong").Text()
	count, err := strconv.ParseInt(number, 10, 0)
	if err != nil && number != "keine " {
		fmt.Println(err)
	}
	return int(count)
}

func getTimestamp(s *goquery.Selection) time.Time {
	ts, err := time.Parse("(Stand: 02.01.2006, 15:04 Uhr)", s.Text())
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return ts

}

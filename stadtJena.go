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

const defaultStatus = "Für Jena liegen aktuell keine Meldungen über Erkrankungen vor."

func getStatus(s *goquery.Selection) string {
	html, err := s.Children().First().Html()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.Trim(textize(html), " ")
}

func getCount(s *goquery.Selection) int64 {
	number := s.Find("strong").Text()
	count, err := strconv.ParseInt(number, 10, 0)
	if err != nil && number != "keine " {
		fmt.Println(err)
	}
	return count
}

func getTimestamp(s *goquery.Selection) time.Time {
	ts, err := time.Parse("(Stand: 02.01.2006, 15:04 Uhr)", s.Text())
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return ts

}

func stadtNeuigkeiten(lastUpdate, lastCount *int64, h *colly.HTMLElement) {
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

	if *lastCount == 0 && statusInfo != defaultStatus {
		sendSignal("%s WARNUNG - CORONA FALL\n%s", zeitpunkt, statusInfo)
		*lastCount = 1
	} else if count > *lastCount {
		sendSignal("%s Stadt Statusänderung\nGestiegen von %d auf %d\n%s", zeitpunkt, lastCount, count, statusInfo)
		*lastCount = count
	} else if timestamp.Unix() > *lastUpdate {
		sendSignal("%s Stadt Aktualisierung\n%s", zeitpunkt, statusInfo)
		*lastUpdate = timestamp.Unix()
	}
}

func stadtJena(lastUpdate, lastCount *int64) {
	c := colly.NewCollector()
	c.OnHTML(".content-inner--main > div > div:first-child > .paragraph--type--text > .text-formatted:first-child", func(h *colly.HTMLElement) { stadtNeuigkeiten(lastUpdate, lastCount, h) })
	c.Visit("http://jena.de/corona")
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/goodsign/monday"
)

var timestamp time.Time
var newest int64
var text string
var entries []string
var level = 0

func ladeBlog(i int, e *goquery.Selection) bool {
	// fmt.Println(i, level, goquery.NodeName(e))
	switch goquery.NodeName(e) {
	case "h2":
		switch level {
		case 3:
			if timestamp.Unix() <= s.timestamp {
				return false
			}
			if timestamp.Unix() > newest {
				newest = timestamp.Unix()
			}
			zeitpunkt := monday.Format(timestamp, "2. January 15:04 Uhr", "de_DE")
			// fmt.Printf("%s\n%s\n\n", zeitpunkt, text)
			entries = append(entries, fmt.Sprintf("%s\n%s", zeitpunkt, text))
			text = ""
			level = 0
			fallthrough
		case 0:
			tt, err := monday.Parse("2. January", e.Text(), "de_DE")
			if err == nil {
				level++
				timestamp = tt
				return true
			}
			fallthrough
		case 1, 2:
			ueberschrift := strings.SplitN(e.Text(), "Uhr", 2)
			zeit := strings.Trim(ueberschrift[0], " ")
			tt, err := time.Parse("15", zeit)
			if err != nil {
				tt, err = time.Parse("15.04", zeit)
			}
			if err != nil {
				tt, err = time.Parse("15:04", zeit)
			}
			if err == nil {
				level = 2
				if len(ueberschrift) > 1 {
					text = strings.Trim(ueberschrift[1], " -:")
					if len(text) > 0 {
						level++
					}
				}
				timestamp = time.Date(
					2020, timestamp.Month(), timestamp.Day(),
					tt.Hour(), tt.Minute(), 0, 0,
					timestamp.Location())
				return true
			}
			if level == 1 {
				text = zeit
				level++
				return true
			}
			if i == 0 {
				return true
			}
			fmt.Println(i, "Kann ich nicht zuordnen", err, ueberschrift, e.Text())
			return true
		default:
			fmt.Println("Surprising level", level)
			return true
		}
	case "p":
		text += "\n" + e.Text()
		return true
	}
	panic("Should never get here")
}

func otzBlog(filename string) {
	s.load(filename)
	newest = s.timestamp
	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find("h2, .article__paragraph").EachWithBreak(ladeBlog)
		for i := range entries {
			sendSignal(entries[len(entries)-i-1])
		}
		s.timestamp = newest
		s.save()
	})
	c.Visit("https://www.otz.de/leben/gesundheit-medizin/coronavirus-thueringen-saale-orla-kreis-verdacht-infiziert-schutzmassnahmen-id228564867.html")
}

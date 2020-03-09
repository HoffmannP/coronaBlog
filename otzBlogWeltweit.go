package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/goodsign/monday"
)

type parseStateWW struct {
	f                bool
	n                string
	newestEntry      int64
	timeCurrentEntry time.Time
	textCurrentEntry string
	entries          []string
}

func ladeBlogWW(lastUpdate *int64, p *parseStateWW, i int, e *goquery.Selection) bool {
	switch goquery.NodeName(e) {
	case "h2":
		ueberschriftDatum := strings.Trim(strings.SplitN(e.Text(), ": ", 2)[0], " ")
		tt, err := monday.Parse("Monday, 2. January 2006", ueberschriftDatum, "de_DE")
		if err == nil {
			p.timeCurrentEntry = tt
		}
		return true
	case "p":
		rawText := strings.Trim(e.Text(), " \n")
		strong := strings.SplitN(e.Find("strong").First().Text(), "Uhr", 2)
		if len(strong) == 1 {
			p.textCurrentEntry += " " + rawText
			return true
		}
		strongText := strings.Trim(strong[0], " \n")
		if 4 > len(strongText) || len(strongText) > 5 {
			p.textCurrentEntry += " " + rawText
			// fmt.Println(len(strongText), "ist keine Uhrzeit")
			return true
		}
		tt, err := time.Parse("15.04", strongText)
		if err != nil {
			// fmt.Println(strongText, "ist doch keine Uhrzeit")
			p.textCurrentEntry += " " + rawText
			return true
		}
		if *lastUpdate >= p.timeCurrentEntry.Unix() {
			// fmt.Println(p.timeCurrentEntry, "ist keine Neuigkeit")
			return false
		}
		if p.timeCurrentEntry.Unix() > p.newestEntry {
			p.newestEntry = p.timeCurrentEntry.Unix()
		}
		zeitpunkt := monday.Format(p.timeCurrentEntry, "2. January 15:04 Uhr", "de_DE")
		if p.f {
			p.entries = append(p.entries, fmt.Sprintf("%s\n%s", zeitpunkt, p.textCurrentEntry))
		} else {
			p.f = true
		}
		p.timeCurrentEntry = time.Date(
			2020, p.timeCurrentEntry.Month(), p.timeCurrentEntry.Day(),
			tt.Hour(), tt.Minute(), 0, 0,
			p.timeCurrentEntry.Location())
		p.textCurrentEntry = strings.Trim(rawText[9:], ".: \n")
		return true
	}
	fmt.Println(goquery.NodeName(e))
	panic("Irgendwas anderes")
}

func otzBlogWeltweit(lastUpdate *int64) {
	today := time.Date(
		2020, time.Now().Month(), time.Now().Day(),
		0, 0, 0, 0,
		time.Now().Location())
	p := parseStateWW{n: "Weltweit", timeCurrentEntry: today}
	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find("h2, .article__paragraph").EachWithBreak(func(i int, e *goquery.Selection) bool { return ladeBlogWW(lastUpdate, &p, i, e) })
		if p.newestEntry > *lastUpdate {
			*lastUpdate = p.newestEntry
		}
		for i := range p.entries {
			sendSignal(p.entries[len(p.entries)-i-1])
		}
	})
	c.Visit("https://www.otz.de/leben/vermischtes/coronavirus-news-italien-norden-abgeriegelt-mailand-und-venedig-betroffen-id228637475.html")
}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/goodsign/monday"
)

type parseStateTH struct {
	n                string
	newestEntry      int64
	timeCurrentEntry time.Time
	textCurrentEntry string
	level            int
	entries          []string
}

func ladeBlogTH(lastUpdate *int64, p *parseStateTH, i int, e *goquery.Selection) bool {
	// fmt.Println(i, p.level, goquery.NodeName(e))
	textContent := strings.Trim(e.Text(), " ")
	switch goquery.NodeName(e) {
	case "h2":
		switch p.level {
		case 3:
			if *lastUpdate >= p.timeCurrentEntry.Unix() {
				return false
			}
			if p.timeCurrentEntry.Unix() > p.newestEntry {
				p.newestEntry = p.timeCurrentEntry.Unix()
			}
			zeitpunkt := monday.Format(p.timeCurrentEntry, "2. January 15:04 Uhr", "de_DE")
			// fmt.Printf("%s\n%s\n\n", zeitpunkt, p.textCurrentEntry)
			p.entries = append(p.entries, fmt.Sprintf("%s\n%s", zeitpunkt, p.textCurrentEntry))
			p.textCurrentEntry = ""
			p.level = 0
			fallthrough
		case 0:
			tt, err := monday.Parse("2. January", textContent, "de_DE")
			if err != nil {
				tt, err = monday.Parse("Monday, 2. January 2006", textContent, "de_DE")
			}
			if err == nil {
				p.level++
				p.timeCurrentEntry = tt
				return true
			}
			fallthrough
		case 1, 2:
			if strings.Index(textContent, "Uhr") == -1 {
				textContent = strings.Replace(textContent, ": ", "Uhr", 1)
			}
			ueberschrift := strings.SplitN(textContent, "Uhr", 2)
			zeit := strings.Trim(ueberschrift[0], " ")
			tt, err := time.Parse("15", zeit)
			if err != nil {
				tt, err = time.Parse("15.04", zeit)
			}
			if err != nil {
				tt, err = time.Parse("15:04", zeit)
			}
			if err == nil {
				p.level = 2
				if len(ueberschrift) > 1 {
					p.textCurrentEntry = strings.Trim(ueberschrift[1], " -:")
					if len(p.textCurrentEntry) > 0 {
						p.level++
					}
				}
				p.timeCurrentEntry = time.Date(
					2020, p.timeCurrentEntry.Month(), p.timeCurrentEntry.Day(),
					tt.Hour(), tt.Minute(), 0, 0,
					p.timeCurrentEntry.Location())
				return true
			}
			if p.level == 1 {
				p.textCurrentEntry = zeit
				p.level++
				return true
			}
			if i == 0 {
				return true
			}
			html, _ := e.Html()
			fmt.Printf("---\nOTZ %s Zeile %d (Level %d) kann nicht zugeordnet werden\n%s\n<h2>%s</h2>\n", p.n, i, p.level, err, html)
			return true
		default:
			fmt.Println("Surprising p.level", p.level)
			return true
		}
	case "p":
		p.textCurrentEntry += "\n" + textContent
		return true
	}
	panic("Should never get here")
}

func otzBlogThueringen(lastUpdate *int64) {
	p := parseStateTH{n: "Thüringen"}
	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.DOM.Find("h2, .article__paragraph").EachWithBreak(func(i int, e *goquery.Selection) bool { return ladeBlogTH(lastUpdate, &p, i, e) })
		if p.newestEntry > *lastUpdate {
			*lastUpdate = p.newestEntry
		}
		for i := range p.entries {
			sendSignal(p.entries[len(p.entries)-i-1])
		}
	})
	c.Visit("https://www.otz.de/leben/gesundheit-medizin/coronavirus-thueringen-saale-orla-kreis-verdacht-infiziert-schutzmassnahmen-id228564867.html")
}

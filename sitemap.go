package main

import (
	"encoding/xml"
	"fmt"
)

type Url struct {
	Localization string `xml:"loc"`
	Priority     string `xml:"priority"`
}

type Sitemap struct {
	XMLName string `xml:"urlset"`
	Urls    []Url  `xml:"url"`
	Xmlns   string `xml:"xmlns,attr"`
}

func NewSitemap(urls []string) *Sitemap {
	u := make([]Url, len(urls))
	for i, url := range urls {
		u[i] = Url{
			Localization: url,
			Priority:     fmt.Sprintf("%.1f", 1.0),
		}
	}
	return &Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  u,
	}
}

func (s *Sitemap) ToXML() ([]byte, error) {
	output, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}

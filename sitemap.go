package main

import "encoding/xml"

type SiteMapXML struct {
	XMLName  xml.Name `xml:"urlset"`
	Version  string   `xml:"xmlns,attr"`
	Xhtml    string   `xml:"xmlns:xhtml,attr"`
	SiteList []*Site  `xml:"url"`
}

type Site struct {
	URL       string `xml:"loc"`
	UpdatedAt string `xml:"changefreq"`
}

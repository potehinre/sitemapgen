package sitemap

import "encoding/xml"

const (
	xmlIdent     = "  "
	sitemapXMLNS = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type Url struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	Lastmod    string   `xml:"lastmod,omitempty"`
	Changefreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"urls"`
}

func Generate(urls []string) []byte {
	urlsSM := []Url{}
	for _, url := range urls {
		urlsSM = append(urlsSM, Url{Loc: url})
	}
	urlset := Urlset{Urls: urlsSM, Xmlns: sitemapXMLNS}
	output, _ := xml.MarshalIndent(urlset, "", xmlIdent)
	res := []byte(xml.Header)
	res = append(res, output...)
	return res
}

package feed

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/mreiferson/go-httpclient"
	"github.com/rogpeppe/go-charset/charset"
	_ "github.com/rogpeppe/go-charset/data" //initialize only
)

const (
	WORDPRESS_DATE_FORMAT               = "Mon, 02 Jan 2006 15:04:05 -0700"
	DEFAULT_TIMEOUT       time.Duration = 10  // seconds
	RESPONSE_TIMEOUT      time.Duration = 120 // seconds
)

//Fetcher interface
type Fetcher interface {
	Get(url string) (resp *http.Response, err error)
}

//Date type
type RSSDate string

//Channel struct for RSS
type Channel struct {
	Title         string      `xml:"title"`
	Link          string      `xml:"link"`
	Description   string      `xml:"description"`
	Language      string      `xml:"language"`
	LastBuildDate RSSDate     `xml:"lastBuildDate"`
	Image         ImageAsset  `xml:"image"`
	Subtitle      string      `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd subtitle"`
	Owner         ItunesOwner `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd owner"`
	Item          []Item      `xml:"item"`
}

//Item struct for each Item in the Channel
type Item struct {
	Title       string          `xml:"title"`
	Link        string          `xml:"link"`
	Comments    string          `xml:"comments"`
	PubDate     RSSDate         `xml:"pubDate"`
	GUID        string          `xml:"guid"`
	Category    []string        `xml:"category"`
	Enclosure   []ItemEnclosure `xml:"enclosure"`
	Description string          `xml:"description"`
	Text        string          `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Duration    string          `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd duration"`
	Author      string          `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd author"`
}

//ItemEnclosure struct for each Item Enclosure
type ItemEnclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type ImageAsset struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

// iTunes namespace
type ItunesOwner struct {
	Name  string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd name"`
	Email string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd email"`
}

//Parse (Date function) and returns Time, error
func (d RSSDate) Parse() (time.Time, error) {
	t, err := d.ParseWithFormat(WORDPRESS_DATE_FORMAT)
	if err != nil {
		t, err = d.ParseWithFormat(time.RFC1123) // variation of the wordpress format
		if err != nil {
			t, err = d.ParseWithFormat(time.RFC822) // RSS 2.0 spec
		}
	}
	return t, err
}

//ParseWithFormat (Date function), takes a string and returns Time, error
func (d RSSDate) ParseWithFormat(format string) (time.Time, error) {
	return time.Parse(format, string(d))
}

//Format (Date function), takes a string and returns string, error
func (d RSSDate) Format(format string) (string, error) {
	t, err := d.Parse()
	if err != nil {
		return "", err
	}
	return t.Format(format), nil
}

//MustFormat (Date function), take a string and returns string
func (d RSSDate) MustFormat(format string) string {
	s, err := d.Format(format)
	if err != nil {
		return err.Error()
	}
	return s
}

//Read a string url and returns a Channel struct, error
func RSS(url string) (*Channel, error) {
	//return ReadWithClient(url, http.DefaultClient)
	transport := &httpclient.Transport{
		ConnectTimeout:        DEFAULT_TIMEOUT * time.Second,
		RequestTimeout:        RESPONSE_TIMEOUT * time.Second,
		ResponseHeaderTimeout: DEFAULT_TIMEOUT * time.Second,
	}
	defer transport.Close()

	client := &http.Client{Transport: transport}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	xmlDecoder := xml.NewDecoder(response.Body)
	xmlDecoder.CharsetReader = charset.NewReader

	var rss struct {
		Channel Channel `xml:"channel"`
	}
	if err = xmlDecoder.Decode(&rss); err != nil {
		return nil, err
	}
	return &rss.Channel, nil
}

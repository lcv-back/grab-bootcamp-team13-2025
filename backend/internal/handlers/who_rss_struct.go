package handlers

type WhoRss struct {
	Channel struct {
		Items []WhoRssItem `xml:"item"`
	} `xml:"channel"`
}

type WhoRssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

package helper

import (
	"github.com/anaskhan96/soup"
	"net/url"
	"strings"
)

type searchResult struct {
	Title string `json:"Title"`
	Path  string `json:"path"`
}

type tags struct {
	Mp3_av  bool `json:"mp3"`
	Flac_av bool `json:"flac"`
	Ogg_av  bool `json:"ogg"`
}

type tracks struct {
	Title  string `json:"Title"`
	Path   string `json:"path"`
	Length string `json:"length"`
}

type TrackResult struct {
	Tags     tags     `json:"tags"`
	AlbumArt string   `json:"albumart"`
	Tracks   []tracks `json:"tracks"`
}

type DirectResult struct {
	Link string `json:"link"`
}

var (
	base_url                  = "https://downloads.khinsider.com/"
	base_search_url           = "search?search="
	base_soundtrack_album_url = "game-soundtracks/album/"
)

func SearchFor(term string) []searchResult {
	var searchResults []searchResult
	var linkArray []string
	var textArray []string

	var completedSearch = base_url + base_search_url + url.QueryEscape(term)

	resp, err := soup.Get(completedSearch)
	if err != nil {
		return []searchResult{
			{"nil", "nil"},
		}
	}

	doc := soup.HTMLParse(resp)
	links := doc.Find("div", "id", "EchoTopic").FindAll("p") //.FindAll("a")

	for _, link := range links[1].FindAll("a") {
		textArray = append(textArray, link.Text())
		linkArray = append(linkArray, link.Attrs()["href"])
	}

	for i := range textArray {
		searchResults = append(searchResults, []searchResult{
			{textArray[i], linkArray[i]},
		}...)
	}
	return searchResults
}

func TrackGetter(downloadCode string) TrackResult {
	var completedTrack = base_url + base_soundtrack_album_url + downloadCode
	var mp3 = false
	var flac = false
	var ogg = false
	var albumart = "none"
	var tracksall []tracks

	resp, err := soup.Get(completedTrack)
	if err != nil {
		return TrackResult{tags{false, false, false}, "nil", tracksall}
	}

	doc := soup.HTMLParse(resp)

	for _, text := range doc.FindAll("p") {
		if text.Text() == "No such album" {
			return TrackResult{tags{false, false, false}, "nil", tracksall}
		}
	}

	var albumartHTML = doc.Find("div", "id", "EchoTopic").FindAll("table")[0] //.Find("img")//.Attrs()["src"]

	for _, imageHTML := range albumartHTML.FindAll("td") {
		if len(imageHTML.Children()) >= 2 {
			albumart = imageHTML.Find("img").Attrs()["src"]
		}
	}

	songlist := doc.Find("table", "id", "songlist")
	header := songlist.Find("tr", "id", "songlist_header")

	var songLoc = 2
	for i := range header.FindAll("th") {
		if header.FindAll("th")[i].Attrs()["width"] == "60px" {
			switch header.FindAll("th")[i].Find("b").Text() {
			case "MP3":
				{
					mp3 = true
				}
			case "FLAC":
				{
					flac = true
				}
			case "OGG":
				{
					ogg = true
				}
			}
		} else if header.FindAll("th")[i].Attrs()["colspan"] != "" {
			songLoc = i
		}
	}

	for _, row := range songlist.FindAll("tr") {
		if row.Attrs()["id"] == "songlist_header" || row.Attrs()["id"] == "songlist_footer" {
			continue
		}

		tracksall = append(tracksall, tracks{row.FindAll("td")[songLoc].Find("a").Text(),
			strings.Replace(row.FindAll("td")[songLoc].Find("a").Attrs()["href"], "/"+base_soundtrack_album_url, "", 1),
			row.FindAll("td")[songLoc+1].Find("a").Text()})
	}

	return TrackResult{
		tags{Mp3_av: mp3, Ogg_av: ogg, Flac_av: flac},
		albumart,
		tracksall,
	}
}

func DirectLinkGrab(downloadCode string, trackName string, format string) DirectResult {
	var completeURL = base_url + base_soundtrack_album_url + downloadCode + "/" + trackName
	resp, err := soup.Get(completeURL)
	if err != nil {
		return DirectResult{"nil"}
	}

	doc := soup.HTMLParse(resp)
	links := doc.Find("div", "id", "EchoTopic")

	for _, tags := range links.FindAll("a") {
		if tags.Attrs()["href"] != "" {
			if strings.HasSuffix(tags.Attrs()["href"], format) {
				return DirectResult{tags.Attrs()["href"]}
			}
		}
	}
	return DirectResult{"nil"}
}

package commands

import (
	"encoding/xml"
	"errors"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/russross/blackfriday"
)

type appXML struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Sparkle string     `xml:"xmlns:sparkle,attr"`
	DC      string     `xml:"xmlns:dc,attr"`
	Channel channelXML `xml:"channel"`
}

type channelXML struct {
	Title string  `xml:"title"`
	Link  string  `xml:"link"`
	Item  itemXML `xml:"item"`
}

type itemXML struct {
	Title       string       `xml:"title"`
	Description cdataType    `xml:"description"`
	Date        string       `xml:"pubDate"`
	Enclosure   enclosureXML `xml:"enclosure"`
}
type cdataType struct {
	Chardata string `xml:",cdata"`
}

type enclosureXML struct {
	URL          string `xml:"url,attr"`
	Version      string `xml:"sparkle:version,attr"`
	ShortVersion string `xml:"sparkle:shortVersionString,attr"`
	Length       string `xml:"length,attr"`
	Type         string `xml:"type,attr"`
}

func genXML(release appXML, tmpRelease string, description []byte) error {
	desc := string(blackfriday.MarkdownCommon(description))
	release.Channel.Item.Description = cdataType{Chardata: desc}
	xmlByte, err := xml.Marshal(release)
	if err != nil {
		logrus.Error("Could not generate XML")
		return errors.New("Could not gen XML")
	}
	logrus.Debug("XML: ", string(xmlByte))
	err = ioutil.WriteFile(tmpRelease, xmlByte, 0644)
	if err != nil {
		logrus.Error("Could not write XML to tmp")
		return errors.New("Could not write XML in tmp: " + tmpRelease)
	}
	return nil
}

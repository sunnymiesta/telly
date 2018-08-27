// Package guide_providers is a telly internal package to provide electronic program guide (EPG) data.
// It is generally modeled after the XMLTV standard with slight deviations to accomodate other providers.
package guide_providers

import (
	"strings"

	"github.com/tellytv/telly/internal/xmltv"
)

type Configuration struct {
	Name     string `json:"-"`
	Provider string

	// Only used for Schedules Direct provider
	Username string
	Password string
	Lineups  []string

	// Only used for XMLTV provider
	XMLTVURL string
}

func (i *Configuration) GetProvider() (GuideProvider, error) {
	switch strings.ToLower(i.Provider) {
	case "schedulesdirect", "schedules-direct", "sd":
		return newSchedulesDirect(i)
	default:
		return newXMLTV(i)
	}
}

// Channel describes a channel available in the providers lineup with necessary pieces parsed into fields.
type Channel struct {
	// Required Fields
	ID     string
	Name   string
	Logos  []Logo
	Number string

	// Optional fields
	CallSign string
	URLs     []string
	Lineup   string
}

func (c *Channel) XMLTV() xmltv.Channel {
	ch := xmltv.Channel{
		ID:   c.ID,
		LCN:  c.Number,
		URLs: c.URLs,
	}

	// Why do we do this? From tv_grab_zz_sdjson:
	//
	// MythTV seems to assume that the first three display-name elements are
	// name, callsign and channel number. We follow that scheme here.
	ch.DisplayNames = []xmltv.CommonElement{
		xmltv.CommonElement{
			Value: c.Name,
		},
		xmltv.CommonElement{
			Value: c.CallSign,
		},
		xmltv.CommonElement{
			Value: c.Number,
		},
	}

	for _, logo := range c.Logos {
		ch.Icons = append(ch.Icons, xmltv.Icon{
			Source: logo.URL,
			Width:  logo.Width,
			Height: logo.Height,
		})
	}

	return ch
}

// A Logo stores the information about a channel logo
type Logo struct {
	URL    string `json:"URL"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// GuideProvider describes a IPTV provider configuration.
type GuideProvider interface {
	Name() string
	Channels() ([]Channel, error)
	Schedule(channelIDs []string) ([]xmltv.Programme, error)

	Refresh() error
	Configuration() Configuration
}

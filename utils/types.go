package utils

import "strings"

type Client struct {
	Name    string `json:"clientName,omitempty"`
	Version string `json:"clientVersion,omitempty"`
	Hl      string `json:"hl",omitempty`
}

type Context struct {
	Client Client `json:"client",omitempty`
}

type PageType struct {
	PageType string `json:"pageType,omitempty"`
}

type BrowseMusicConfig struct {
	MusicConfig PageType `json:"browseEndpointContextMusicConfig,omitempty"`
}

type Form struct {
	Values []string `json:"selectedValues"`
}

type WatchMusicConfig struct {
	Panel bool   `json:"hasPersistentPlaylistPanel,omitempty"`
	Type  string `json:"musicVideoType,omitempty"`
}

type BrowseData struct {
	Context     Context           `json:"context,omitempty"`
	MusicConfig BrowseMusicConfig `json:"browseEndpointContextMusicConfig,omitempty"`
	BrowseId    string            `json:"browseId,omitempty"`
	Params      string            `json:"params,omitempty"`
	Form        Form              `json:"formData,omitempty"`
}

type NextData struct {
	Id          string           `json:"videoId,omitempty"`
	PlaylistId  string           `json:"playlistId",omitempty`
	Context     Context          `json:"context,omitempty"`
	Audio       bool             `json:"isAudioOnly,omitempty"`
	Tuner       string           `json:"tunerSettingValue,omitempty"`
	Panel       bool             `json:"enablePersistentPlaylistPanel,omitempty"`
	MusicConfig WatchMusicConfig `json:"watchEndpointMusicConfig,omitempty"`
}

var BaseContext Context = Context{
	Client: Client{
		Name:    "WEB_REMIX",
		Version: "1.20220926.01.00",
		Hl:      "en-US",
	},
}

func TypeBrowse(t, id, params, form string) BrowseData {
	if t != "" {
		return BrowseData{
			Context: BaseContext,
			MusicConfig: BrowseMusicConfig{
				MusicConfig: PageType{
					PageType: "MUSIC_PAGE_TYPE_" + strings.ToUpper(t),
				},
			},
			BrowseId: id,
		}
	} else if form != "" {
		return BrowseData{
			Context:  BaseContext,
			BrowseId: id,
			Params:   params,
			Form: Form{
				Values: []string{form},
			},
		}
	} else {
		return BrowseData{
			Context:  BaseContext,
			BrowseId: id,
			Params:   params,
		}
	}
}

func TypeNext(id, pid string) NextData {
	return NextData{
		Id:         id,
		PlaylistId: pid,
		Context:    BaseContext,
		Panel:      true,
		Audio:      true,
		Tuner:      "AUTOMIX_SETTING_NORMAL",
		MusicConfig: WatchMusicConfig{
			Panel: true,
			Type:  "MUSIC_VIDEO_TYPE_ATV",
		},
	}
}

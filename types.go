package main

import "strings"

var BaseContext Context = Context{
	Client: Client{
		Name:    "WEB_REMIX",
		Version: "1.20220926.01.00",
		Hl:      "en-US",
	},
}

func GetTypeBrowse(t, id, params, form string) BrowseData {
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

func GetTypeNext(id, pid string) NextData {
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

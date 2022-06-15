package main

import "strings"

var BaseContext = Context{
	Client: Client{
		Name:    "WEB_REMIX",
		Version: "1.20211213.00.00",
	},
}

func GetTypeBrowse(t, id string) BrowseData {

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
	} else {
		return BrowseData{
			Context:  BaseContext,
			BrowseId: id,
		}
	}

}

func GetTypeNext(id string) NextData {
	return NextData{
		Id:      id,
		Context: BaseContext,
		Panel:   true,
		Audio:   true,
		Tuner:   "AUTOMIX_SETTING_NORMAL",
		MusicConfig: WatchMusicConfig{
			Panel: true,
			Type:  "MUSIC_VIDEO_TYPE_ATV",
		},
	}
}

package main

import "strings"

func GetContext() Context {
	return Context{
		Client: Client{
			Name:    "WEB_REMIX",
			Version: "1.20211213.00.00",
		},
	}
}

func GetTypeBrowse(t, id string) BrowseData {
	return BrowseData{
		Context: GetContext(),
		MusicConfig: MusicConfig{
			MusicConfig: PageType{
				PageType: "MUSIC_PAGE_TYPE_" + strings.ToUpper(t),
			},
		},
		BrowseId: id,
	}
}

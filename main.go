package main

import (
	"encoding/json"
	"fmt"
	"github.com/hjson/hjson-go"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/shibukawa/configdir"
	"io/ioutil"
)

const (
	width         = 400
	height        = 191
	defaultVolume = 70
)

func load(file string) *Config {
	b, e := ioutil.ReadFile(file)
	if e != nil {
		panic(e)
	}
	c := Config{}
	var dat map[string]interface{}
	e = hjson.Unmarshal(b, &dat)
	if e != nil {
		panic(e)
	}
	// convert to JSON
	jj, _ := json.Marshal(dat)
	// unmarshal
	json.Unmarshal(jj, &c)

	return &c

}

func main() {
	configDirs := configdir.New("osjuptier", "GMRepeater")
	cache := configDirs.QueryCacheFolder()

	firstVolume := defaultVolume
	if cache.Exists("setting.json") {
		var data PlayConf
		str, _ := cache.ReadFile("setting.json")
		if e := json.Unmarshal(str, &data); e != nil {
			panic(e)
		}
		firstVolume = data.Volume
	}

	conf := load("config.hjson")

	player := NewPlayer(conf)

	var model walk.ListModel = &SongListItems{
		songList: player.config.Songs,
	}

	var volumeSlide *walk.Slider
	var songName *walk.Label
	var nowPlaying *walk.Label
	var combo *walk.ComboBox

	player.callback = func(id int, now float64, all float64) {
		songName.SetText(player.config.Songs[id].Name)
		n := int(now)
		m := int(all)
		nowM, nowS := n/60, n%60
		maxM, maxS := m/60, m%60
		nowPlaying.SetText(fmt.Sprintf("%02d:%02d / %02d:%02d ", nowM, nowS, maxM, maxS))
	}
	icon, iconErr := walk.Resources.Icon("4")
	if iconErr != nil {
		panic(iconErr)
	}

	size := Size{Width: width, Height: height}
	var main *walk.MainWindow
	e := MainWindow{
		AssignTo: &main,
		Icon:     icon,
		Title:    "GMRepeater",
		MaxSize:  size,
		MinSize:  size,
		Size:     size,
		Layout:   VBox{},
		Children: []Widget{
			ComboBox{
				AssignTo:      &combo,
				Value:         Bind("SpeciesId", SelRequired{}),
				BindingMember: "Id",
				DisplayMember: "Name",
				Model:         model,
			},
			Composite{
				Layout:    VBox{},
				Alignment: AlignHCenterVCenter,
				Children: []Widget{
					Label{
						AssignTo:  &songName,
						Text:      " - ",
						Alignment: AlignHCenterVCenter,
					},
					Label{
						AssignTo:  &nowPlaying,
						Text:      "00:00 / 00:00",
						Alignment: AlignHCenterVCenter,
					},
				},
			},
			Composite{
				Layout: Flow{},
				Children: []Widget{
					PushButton{
						Text: "Start",
						OnClicked: func() {
							id := combo.CurrentIndex()
							player.Start(id)
						},
					},
					PushButton{
						Text: "Stop",
						OnClicked: func() {
							player.Suspend()
						},
					},
					Slider{
						AssignTo: &volumeSlide,
						Value:    firstVolume,
						MaxValue: 100,
						MinValue: 0,
						Tracking: true,
						OnValueChanged: func() {
							player.SetVol(volumeSlide.Value())

						},
					},
				},
			},
		},
		OnBoundsChanged: func() {
			main.SetSize(walk.Size{Width: width, Height: height})
		},
		OnSizeChanged: func() {
			main.SetSize(walk.Size{Width: width, Height: height})
		},
	}.Create()

	main.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {

		data := PlayConf{Volume: volumeSlide.Value()}
		str, _ := json.Marshal(&data)
		if e := cache.WriteFile("setting.json", str); e != nil {
			panic(e)
		}

	})
	if e != nil {
		panic(e)
	}

	main.Run()

}

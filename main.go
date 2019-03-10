package main

import (
	"encoding/json"
	"fmt"
	"github.com/hjson/hjson-go"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
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
	conf := load("config.hjson")

	player := NewPlayer(conf)

	var model walk.ListModel = &SongListItems{
		songList: player.config.Songs,
	}

	var volumeSlide *walk.Slider
	var songName *walk.Label
	var nowPlaying *walk.Label
	var combo *walk.ComboBox

	player.callback = func(now float64, all float64) {
		songName.SetText(player.config.Songs[0].Name)
		n := int(now)
		m := int(all)
		nowM, nowS := n/60, n%60
		maxM, maxS := m/60, m%60
		nowPlaying.SetText(fmt.Sprintf("%02d:%02d / %02d:%02d ", nowM, nowS, maxM, maxS))
	}

	_, e := MainWindow{
		Title:  "OGG Repeat",
		Size:   Size{Width: 400, Height: 100},
		Layout: VBox{},
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
						Value:    100,
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
	}.Run()
	if e != nil {
		panic(e)
	}

}

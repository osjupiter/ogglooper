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
	var playingSlide *walk.Slider
	var nowPlaying *walk.Label
	var combo *walk.ComboBox

	player.secondHandler = func(now float64,all float64) {
		percent := now / (all) * 100
		playingSlide.SetValue(int(percent))
		songName.SetText(player.config.Songs[0].File)

		n:=int(now)
		m:=int(all)
		nowM,nowS:=n/60,n%60
		maxM,maxS:=m/60,m%60
		nowPlaying.SetText(fmt.Sprintf("%02d:%02d / %02d:%02d ",nowM,nowS ,maxM,maxS))
	}

	// ドロップダウン
	// 再生停止
	// 一時停止
	// ループする
	// 音量

	// Opt 一時間
	// 再生位置
	// 閉じるボタン（フェードアウト？
	// 最小化
	_, e := MainWindow{
		Title:   "OGG Repeat",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			Slider{
				AssignTo: &playingSlide,
			},
			Label{
				AssignTo: &songName,
				Text:     "なにか",
			},
			Label{
				AssignTo: &nowPlaying,
				Text:     "00:00 / 00:00",
			},
			ComboBox{
				AssignTo:      &combo,
				Value:         Bind("SpeciesId", SelRequired{}),
				BindingMember: "Id",
				DisplayMember: "Name",
				Model:         model,
			},
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
	}.Run()
	if e != nil {
		panic(e)
	}

}

type EnvItem struct {
	name  string
	value string
}

type EnvModel struct {
	walk.ListModelBase
	items []EnvItem
}

func (m *EnvModel) ItemCount() int {
	return len(m.items)
}

func (m *EnvModel) Value(index int) interface{} {
	return m.items[index].name
}

package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	player := NewPlayer()

	var model walk.ListModel = &EnvModel{
		items: []EnvItem{
			{name: "1", value: "b"},
			{name: "2", value: "b"},
			{name: "3", value: "b"},
		},
	}

	var slv *walk.Slider
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
			ComboBox{
				Value:         Bind("SpeciesId", SelRequired{}),
				BindingMember: "Id",
				DisplayMember: "Name",
				Model:         model,
			},
			PushButton{
				Text: "Start",
				OnClicked: func() {
					player.Start()
				},
			},
			PushButton{
				Text: "Stop",
				OnClicked: func() {
					player.Suspend()
				},
			},
			Slider{
				AssignTo: &slv,
				Value:    100,
				OnValueChanged: func() {
					player.SetVol(slv.Value())

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

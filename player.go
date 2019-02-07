package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"os"
	"time"
)

type Player struct {
	context   *PlayContext
	list      SongList
	volMaster *effects.Volume
}

func NewPlayer() *Player {
	return &Player{
		context:   &PlayContext{},
		volMaster: &effects.Volume{Base: 10, Volume: 1},
	}
}

func (p *Player) load(list SongList) {
	p.list = list
	// なんか初期化

}
func (p *Player) Start() {

	f, _ := os.Open("1.ogg")
	s, format, _ := vorbis.Decode(f)

	f2, _ := os.Open("2.ogg")
	l, format, _ := vorbis.Decode(f2)
	// import "github.com/faiface/beep/speaker"
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	go p.loop(s, l)
}

func (p *Player) loop(s beep.StreamSeekCloser, l beep.StreamSeekCloser) {
	p.volMaster.Streamer = beep.Seq(s, beep.Loop(-1, l))
	s.Seek(0)
	speaker.Play(p.volMaster)
}
func (p *Player) Suspend() {
	speaker.Lock()

}

/*
func (p *Player) SkipToNext() {

}
func (p *Player) BackToHead() {

}
func (p *Player) PassToNext() {

}
*/
func (p *Player) SetVol(v int) {
	p.volMaster.Volume = 1.0*(float64(v)/100.0) - 1

}
func (p *Player) SelectSong() {

}

/*
- 再生ボタン
- 停止ボタン
- 次にとばすボタン
- 最初に戻る
- この再生で終了するボタン
- 音量調整
- 曲を選択して再生
*/

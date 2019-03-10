package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"time"
)

type Player struct {
	volMaster *effects.Volume
	config    *Config
	callback  func(float64, float64)
	now       beep.StreamCloser
}

func NewPlayer(config *Config) *Player {
	p := &Player{
		volMaster: &effects.Volume{Base: 2, Volume: 0},
	}
	p.config = config
	return p
}

func (p *Player) Start(id int) {
	if id == -1 {
		return
	}

	stream, format := p.config.Songs[id].stream()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	go p.play(stream, format)
}
func (p *Player) doCallback(s beep.StreamSeekCloser, format beep.Format) {
	now := format.SampleRate.D(s.Position()).Round(time.Second)
	max := format.SampleRate.D(s.Len()).Round(time.Second)
	p.callback(now.Seconds(), max.Seconds())
}

func (p *Player) play(s beep.StreamSeekCloser, format beep.Format) {
	if p.now != nil {
		p.now.Close()
	}
	p.volMaster.Streamer = s
	p.doCallback(s, format)

	speaker.Play(p.volMaster)
	p.now = s
	go func() {
		for p.now == s {
			select {
			case <-time.After(time.Second):
				//speaker.Lock()
				p.doCallback(s, format)
				//speaker.Unlock()
			}
		}
	}()
}
func (p *Player) Suspend() {
	speaker.Clear()
}

func (p *Player) SetVol(v int) {
	if v == 0 {
		p.volMaster.Silent = true
	} else {
		p.volMaster.Silent = false
		p.volMaster.Volume = (float64(v)/100.0 - 1) * 10
		fmt.Println(p.volMaster.Volume)
	}
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

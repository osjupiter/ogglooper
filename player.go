package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"os"
	"time"
)

type Player struct {
	context       *PlayContext
	volMaster     *effects.Volume
	config        *Config
	playing       bool
	main          beep.StreamSeekCloser
	intro         beep.StreamSeekCloser
	secondHandler func(float64,float64)
}

func NewPlayer(config *Config) *Player {
	p := &Player{
		context:   &PlayContext{},
		volMaster: &effects.Volume{Base: 10, Volume: 1},
	}
	p.config = config
	return p
}

func (p *Player) Start(id int) {
	if id == -1 {
		return
	}
	main, intro, format := p.stream(id)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	stream := beep.Loop(-1, main)
	if intro != nil {
		stream = beep.Seq(intro, stream)
	}
	p.main = main
	p.intro = intro
	go p.play(stream, format)
}
func (p *Player) stream(id int) (beep.StreamSeekCloser, beep.StreamSeekCloser, beep.Format) {
	songConf := p.config.Songs[id]
	var main beep.StreamSeekCloser = nil
	var mainFormat beep.Format
	var intro beep.StreamSeekCloser = nil
	var introFormat beep.Format

	if songConf.File != "" {
		f2, _ := os.Open(songConf.File)
		main, mainFormat, _ = vorbis.Decode(f2)
	}
	if songConf.IntroFile != "" {
		f, _ := os.Open(songConf.IntroFile)
		intro, introFormat, _ = vorbis.Decode(f)
	}
	if songConf.IntroFile != "" && mainFormat != introFormat {
		panic("err formats")
	}
	return main, intro, mainFormat
}

func (p *Player) play(s beep.Streamer, format beep.Format) {
	p.volMaster.Streamer = s
	speaker.Play(p.volMaster)
	p.playing = true
	go func(playing *bool) {
		for *playing {
			select {
			case <-time.After(time.Second):
				//speaker.Lock()
				now := format.SampleRate.D(p.main.Position() + p.intro.Position()).Round(time.Second)
				max := format.SampleRate.D(p.main.Len() + p.intro.Len()).Round(time.Second)

				p.secondHandler(now.Seconds(), max.Seconds())
				fmt.Println(now)
				//speaker.Unlock()
			}
		}
	}(&p.playing)
}
func (p *Player) Suspend() {
	speaker.Clear()
	p.playing = false
}

func (p *Player) SetVol(v int) {
	p.volMaster.Volume = 3.5*(float64(v)/100.0) - 3.5
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

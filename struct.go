package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/lxn/walk"
	"os"
	"strings"
)

// static
type Config struct {
	Songs           SongList
	DefaultLoopTime int
}
type SongList []*Song

type SongListItems struct {
	songList SongList
	walk.ListModelBase
}

func (s *SongListItems) ItemCount() int {
	return len(s.songList)
}

func (s *SongListItems) Value(index int) interface{} {
	return s.songList[index].Name
}

type Song struct {
	//setting
	Name      string
	File      string
	IntroFile string

	// streamer
	main  beep.StreamSeekCloser
	intro beep.StreamSeekCloser
	loop  beep.Streamer
}

func (s *Song) Stream(samples [][2]float64) (n int, ok bool) {
	return s.loop.Stream(samples)
}
func (s *Song) Err() error {
	return s.loop.Err()
}
func (s *Song) Len() int {
	size := 0
	if s.main != nil {
		size += s.main.Len()
	}
	if s.intro != nil {
		size += s.intro.Len()
	}
	return size
}
func (s *Song) Position() int {
	size := 0
	if s.main != nil {
		size += s.main.Position()
	}
	if s.intro != nil {
		size += s.intro.Position()
	}
	return size

}
func (s *Song) Seek(p int) error {
	panic("not implemented")

}
func (s *Song) Close() error {
	if s.main != nil {
		s.main.Close()
	}
	if s.intro != nil {
		s.intro.Close()
	}
	return nil
}

func readFileAsStreamer(name string) (s beep.StreamSeekCloser, format beep.Format) {
	f, _ := os.Open(name)
	if strings.HasSuffix(name, ".ogg") {
		streamer, format, _ := vorbis.Decode(f)
		return streamer, format
	}
	if strings.HasSuffix(name, ".mp3") {
		streamer, format, _ := mp3.Decode(f)
		return streamer, format
	}
	panic("cant decode the file")

}

func (s *Song) stream() (beep.StreamSeekCloser, beep.Format) {
	songConf := s
	var main beep.StreamSeekCloser = nil
	var mainFormat beep.Format
	var intro beep.StreamSeekCloser = nil
	var introFormat beep.Format

	if songConf.File != "" {
		main, mainFormat = readFileAsStreamer(songConf.File)
	}
	if songConf.IntroFile != "" {
		intro, introFormat = readFileAsStreamer(songConf.IntroFile)
	}
	if songConf.IntroFile != "" && mainFormat != introFormat {
		panic("err formats")
	}

	stream := beep.Loop(-1, main)
	if intro != nil {
		stream = beep.Seq(intro, stream)
	}
	s.main = main
	s.intro = intro
	s.loop = stream
	return s, mainFormat
}

func (s *Song) String() string {
	return fmt.Sprintf("{%s %s %s}", s.Name, s.File, s.IntroFile)
}

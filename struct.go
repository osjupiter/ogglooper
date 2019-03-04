package main

import (
	"bytes"
	"fmt"
	"github.com/lxn/walk"
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

const (
	NONE      State = 0
	PLAYING   State = iota
	SUSPENDED State = iota
)

type Song struct {
	Name      string
	File      string
	IntroFile string
	Data      SongData
}

func (s *Song) String() string {
	return fmt.Sprintf("{%s %s %s}", s.Name, s.File, s.IntroFile)
}

type SongData bytes.Buffer

// dynamic
type PlaySetting struct {
	loopTime int
}
type State int

type PlayContext struct {
	state   State
	playing *Song
	pos     int64
	volume  int
}

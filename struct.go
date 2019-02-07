package main

import "bytes"

// static
type Config struct {
	songs           SongList
	defaultLoopTime int
}
type SongList []*Song

const (
	NONE      State = 0
	PLAYING   State = iota
	SUSPENDED State = iota
)

type Song struct {
	title       string
	RepeatStart int64
	RepeatEnd   int64
	data        SongData
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

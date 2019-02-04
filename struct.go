package main

import "bytes"

// static
type Config struct {
	songs SongList
	defaultLoopTime int
}
type SongList []*Song

type Song struct {
	title       string
	RepeatStart int64
	RepeatEnd   int64
	data SongData
}
type SongData bytes.Buffer

// dynamic
type PlaySetting struct{
	loopTime int
}
type PlayContext struct{
	playing *Song
	pos int64
	volume int
}

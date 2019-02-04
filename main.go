package main

import (
	"os"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/speaker"
	"time"
)

func main() {

	f, _ := os.Open("test.ogg")
	s, format, _ := vorbis.Decode(f)
	// import "github.com/faiface/beep/speaker"
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(s)
	select {} // for now

}

/*

	bu:=bytes.NewBuffer(b)



	r, err := oggvorbis.NewReader(bu)
	if(err!=nil){
		panic(err)
	}
	// handle error

	fmt.Println(r.SampleRate())
	fmt.Println(r.Channels())

	p, err := oto.NewPlayer(44100, 2, 2, 4096)

	buffer := make([]float32, 8192)
	for {
		n, err := r.Read(buffer)

		// use buffer[:n]
		fmt.Println(n)
		fmt.Println(buffer)
		p.Write(buffer)

		if err == io.EOF {
			break
		}
		if err != nil {
			// handle error
		}
	}
	if err := p.Close(); err != nil {
		panic("")
	}


*/

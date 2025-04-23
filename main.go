package main

import (
	"embed"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/inancgumus/screen"
	"log"
	"os"
	"time"
)

type Node[T any] struct {
	data T
	next *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
}

//go:embed res
var res embed.FS

const count = 6569

func main() {
	file, fileError := res.Open("res/BA.wav")
	if fileError != nil {
		log.Fatal(fileError)
	}

	streamer, format, decodeError := wav.Decode(file)
	if decodeError != nil {
		log.Fatal(decodeError)
	}
	defer streamer.Close()
	speakerError := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if speakerError != nil {
		log.Fatal(speakerError)
	}

	frames := &LinkedList[string]{
		head: &Node[string]{},
	}
	index := 1
	head := frames.head
	for index <= count {
		name := fmt.Sprintf("res/BA%d.txt", index)
		bytes, readFileError := res.ReadFile(name)
		if readFileError != nil {
			log.Fatal(readFileError)
		}

		head.data = string(bytes)
		head.next = &Node[string]{}
		head = head.next
		index++
	}

	playing := true

	go func() {
		headLocal := frames.head
		for playing && "" != headLocal.data {
			screen.Clear()
			_, _ = os.Stdout.WriteString(headLocal.data)
			time.Sleep(time.Millisecond * 33) // roughly 30 fps
			headLocal = headLocal.next
		}
	}()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	playing = false
}

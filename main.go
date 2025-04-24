package main

import (
	"embed"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/inancgumus/screen"
	"io"
	"log"
	"os"
	"time"
)

//go:embed res
var res embed.FS

const count = 6569

type Node[T any] struct {
	data T
	next *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
}

func play(playing chan bool) {
	file, openError := res.Open("res/BA.mp3")
	if openError != nil {
		log.Fatal(openError)
	}
	defer file.Close()

	decoder, decoderError := mp3.NewDecoder(file)
	if decoderError != nil {
		log.Fatal(decoderError)
	}

	context, contextError := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if contextError != nil {
		log.Fatal(contextError)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()
	playing <- true
	_, copyError := io.Copy(player, decoder)
	if copyError != nil {
		playing <- false
		log.Fatal(copyError)
	}
	playing <- false
}

func findFrames() *LinkedList[string] {
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
	return frames
}

func draw(playing chan bool) {
	frames := findFrames()
	headLocal := frames.head
	playingAudio := <-playing

	go func() {
		for {
			playingAudio = <-playing
		}
	}()

	for playingAudio && "" != headLocal.data {
		screen.Clear()
		_, _ = os.Stdout.WriteString(headLocal.data)
		time.Sleep(time.Millisecond * 33) // roughly 30 fps
		headLocal = headLocal.next
	}
}

func main() {
	playing := make(chan bool, 1)
	go draw(playing)
	play(playing)
}

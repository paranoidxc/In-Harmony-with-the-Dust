package ui2d

import (
	"github.com/veandco/go-sdl2/mix"
	"math/rand"
	"time"
)

type sounds struct {
	openingDoors []*mix.Chunk
	footsteps    []*mix.Chunk
}

func playRandomSound(chunks []*mix.Chunk, volume int) {
	rand.Seed(int64(time.Now().Unix()))
	chunkIndex := rand.Intn(len(chunks))
	// volume too small will not hear sound
	//fmt.Println("index:", len(chunks), chunkIndex)
	chunks[chunkIndex].Volume(volume)
	chunks[chunkIndex].Play(-1, 0)
}

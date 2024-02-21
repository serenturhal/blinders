package match

import (
	"math/rand"
	"time"
)

// Embedder uses embedding model to embed user information and return embedding vector.
type Embedder interface {
	Embed(info UserMatch) ([]float32, error)
}

type MockEmbedder struct{}

func (e MockEmbedder) Embed(_ UserMatch) ([]float32, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	embed := make([]float32, 10)
	for i := range embed {
		embed[i] = r.Float32()
	}
	return embed, nil
}

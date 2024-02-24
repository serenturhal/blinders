package match

import (
	"math/rand"
	"time"
)

// Embedder uses embedding model to embed user information and return embedding vector.
type Embedder interface {
	Embed(info UserMatch) (EmbeddingVector, error)
}

type MockEmbedder struct{}

func (e MockEmbedder) Embed(_ UserMatch) (EmbeddingVector, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var embed [128]float32
	for i := range embed {
		embed[i] = r.Float32()
	}
	return embed, nil
}

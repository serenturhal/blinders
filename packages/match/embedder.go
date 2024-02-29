package match

import (
	"math/rand"
	"time"

	"blinders/packages/db/models"
)

// Embedder uses embedding model to embed user information and return embedding vector.
type Embedder interface {
	Embed(info models.MatchInfo) (EmbeddingVector, error)
}

type MockEmbedder struct{}

func (e MockEmbedder) Embed(_ models.MatchInfo) (EmbeddingVector, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var embed [384]float32
	for i := range embed {
		embed[i] = r.Float32()
	}
	return embed[:], nil
}

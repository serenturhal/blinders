package match

// Embedder uses embedding model to embed user information and return embedding vector.
type Embedder interface {
	Embed(info UserMatch) ([]float32, error)
}

type MockEmbedder struct{}

func (e MockEmbedder) Embed(_ UserMatch) ([]float32, error) {
	return []float32{}, nil
}

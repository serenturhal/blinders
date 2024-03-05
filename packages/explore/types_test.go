package explore_test

import (
	"blinders/packages/explore"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingThenEmbed(t *testing.T) {
	embed := make(explore.EmbeddingVector, 200)
	for i := range embed {
		embed[i] = rand.Float32()
	}
	encoded, err := embed.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, encoded)

	decodedEmbed := make(explore.EmbeddingVector, 200)
	assert.Nil(t, decodedEmbed.UnmarshalBinary(encoded))

	assert.Equal(t, embed, decodedEmbed)
}

func TestDecodingFail(t *testing.T) {
	embed := make(explore.EmbeddingVector, 200)
	for i := range embed {
		embed[i] = rand.Float32()
	}

	encoded, err := embed.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, encoded)

	decodedEmbed := make(explore.EmbeddingVector, 200)

	wrongEncoded := []byte("failed")
	assert.NotNil(t, decodedEmbed.UnmarshalBinary(wrongEncoded))

	assert.Nil(t, decodedEmbed.UnmarshalBinary(encoded))
	assert.Equal(t, embed, decodedEmbed)
}

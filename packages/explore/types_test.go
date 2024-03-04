package explore

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingEmbed(t *testing.T) {
	embed := make(EmbeddingVector, 200)
	for i := range embed {
		embed[i] = rand.Float32()
	}
	fmt.Println(embed)

	encoded, err := embed.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, encoded)

	decodedEmbed := make(EmbeddingVector, 200)
	assert.Nil(t, decodedEmbed.UnmarshalBinary(encoded))

	assert.Equal(t, embed, decodedEmbed)
}

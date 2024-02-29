package match

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingEmbed(t *testing.T) {
	embed := EmbeddingVector{}
	for i := range embed {
		embed[i] = rand.Float32()
	}

	encoded, err := embed.MarshalBinary()
	assert.Nil(t, err)
	assert.NotNil(t, encoded)
	fmt.Println(len(embed))
	fmt.Println(len(encoded))

	decodedEmbed := new(EmbeddingVector)
	assert.Nil(t, decodedEmbed.UnmarshalBinary(encoded))

	assert.Equal(t, embed, *decodedEmbed)
}

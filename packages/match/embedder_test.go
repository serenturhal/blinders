package match

import (
	"fmt"
	"testing"

	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
)

func TestMockEmbedder(t *testing.T) {
	t.Parallel()
	m := MockEmbedder{}
	numTest := 10000
	for i := 0; i < numTest; i++ {
		t.Run(fmt.Sprintf("test:%v", i), func(t *testing.T) {
			embed, err := m.Embed(models.MatchInfo{})
			assert.Nil(t, err)
			embed2, err := m.Embed(models.MatchInfo{})
			assert.Nil(t, err)
			assert.NotEqual(t, embed, embed2)
		})
	}
}

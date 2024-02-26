package match

import (
	"bytes"
	"encoding/binary"
)

type (
	EmbeddingVector [128]float32
)

func (v EmbeddingVector) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

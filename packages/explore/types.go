package explore

import (
	"bytes"
	"encoding/binary"
)

type (
	EmbeddingVector []float32
)

func (v *EmbeddingVector) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *EmbeddingVector) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	return nil
}

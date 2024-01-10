package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseJSONFile[T any](filename string) (*T, error) {
	bytes, err := GetFile(filename)
	if err != nil {
		return nil, err
	}

	var obj T
	err = json.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, fmt.Errorf("parse json error %v", err)
	}

	return &obj, nil
}

func GetFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("read file error %v", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read content error %v", err)
	}

	return bytes, nil
}

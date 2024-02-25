package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson"
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

func BSONConvert[T any](from any) (*T, error) {
	fromBytes, err := bson.Marshal(from)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}
	var to T
	err = bson.Unmarshal(fromBytes, &to)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}

	return &to, nil
}

func JSONConvert[T any](from any) (*T, error) {
	fromBytes, err := json.Marshal(from)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}
	var to T
	err = json.Unmarshal(fromBytes, &to)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}

	return &to, nil
}

func ParseJSON[T any](from []byte) (*T, error) {
	var to T
	err := json.Unmarshal(from, &to)
	if err != nil {
		return nil, fmt.Errorf("convert error: %v", err)
	}

	return &to, nil
}

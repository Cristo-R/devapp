package utils

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func Struct2Map(i interface{}) (map[string]interface{}, error) {
	var ret map[string]interface{}
	b, err := jsoniter.Marshal(i)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(b, &ret)
	if err != nil {
		return nil, err
	}

	return ret, err
}

func Bytes2Map(bytes []byte) (map[string]interface{}, error) {
	b := jsoniter.NewDecoder(strings.NewReader(string(bytes)))
	b.UseNumber()
	var ret map[string]interface{}
	err := b.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, err
}

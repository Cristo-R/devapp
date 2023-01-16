package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

/*
gorm custom type
*/

type StringArray []string

func (a *StringArray) Scan(src interface{}) error {
	if str, ok := src.([]byte); ok {
		if err := json.Unmarshal([]byte(str), a); err == nil {
			return nil
		}
	}
	return fmt.Errorf("cannot convert %T to StringArray", src)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "[]", nil
	}

	bs, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return string(bs), nil
}

type UUIDBinary string

func (id UUIDBinary) Value() (v driver.Value, err error) {
	if id == "" {
		return nil, nil
	}
	u, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}
	uuidBinary, err := u.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return uuidBinary, nil
}

func (id *UUIDBinary) Scan(value interface{}) (err error) {
	if value == nil {
		*id = ""
		return nil
	}
	s, ok := value.([]byte)

	if !ok {
		*id = ""
		return errors.New("invalid scan source")
	}

	u, err := uuid.FromBytes(s)
	if err != nil {
		*id = ""
		return err
	}
	*id = UUIDBinary(u.String())
	return nil
}

func (id UUIDBinary) String() string {
	return string(id)
}

func (id UUIDBinary) Binary() ([]byte, error) {
	u, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}
	uuidBinary, err := u.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return uuidBinary, nil
}

func (id UUIDBinary) MustBinary() []byte {
	u, err := uuid.Parse(string(id))
	if err != nil {
		panic(err)
	}
	uuidBinary, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return uuidBinary
}

func NewUUIDBinary() UUIDBinary {
	return UUIDBinary(uuid.New().String())
}

var DefaultJSONValue = JSON([]byte("{}"))

type JSON []byte

func NewJSON(a interface{}) (JSON, error) {
	if a == nil {
		return nil, errors.New("NewJSON parameter must not be nil")
	}
	bs, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return JSON(bs), nil
}

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	s, ok := value.([]byte)

	if !ok {
		str, ok := value.(string)
		if ok {
			s = []byte(str)
		} else {
			return errors.New("invalid scan source")
		}
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

func JSONPtr(v JSON) *JSON {
	return &v
}

func JSONBPtr(v []byte) *JSON {
	jsonValue := JSON(v)
	return &jsonValue
}

func (j JSON) UnmarshalToStruct(v interface{}) error {
	if j.IsNull() {
		return errors.New("JSON value is not be null")
	}

	return json.Unmarshal([]byte(j), v)
}

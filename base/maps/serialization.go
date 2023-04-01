package maps

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ToString converts a value to string.
func ToString(value interface{}) string {
	switch value.(type) {
	case string:
		return value.(string)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(value.(int64)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(value.(uint64)), 10)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'g', -1, 64)
	case float64:
		return strconv.FormatFloat(float64(value.(float64)), 'g', -1, 64)
	case bool:
		return strconv.FormatBool(value.(bool))
	default:
		return fmt.Sprintf("%+v", value)
	}
}

// ToJSON outputs the JSON representation of the tree.
func (tree *Map[K, V]) ToJSON() ([]byte, error) {
	elements := make(map[string]interface{})
	it := tree.Iterator()
	for it.Next() {
		elements[ToString(it.Key())] = it.Value()
	}
	return json.Marshal(&elements)
}

// FromJSON populates the tree from the input JSON representation.
func (tree *Map[K, V]) FromJSON(data []byte) error {
	elements := make(map[K]V)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		tree.Clear()
		for key, value := range elements {
			tree.Put(key, value)
		}
	}
	return err
}

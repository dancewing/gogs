package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogits/gogs/pkg/tool"
)

type DataBag map[string]interface{}

func (m *DataBag) Get(keys ...string) (interface{}, bool) {
	return tool.GetMapKey(*m, keys...)
}

func (m *DataBag) GetSlice(keys ...string) ([]interface{}, bool) {
	slice, ok := tool.GetMapKey(*m, keys...)
	if slice != nil {
		return slice.([]interface{}), ok
	}
	return nil, false
}

func (m *DataBag) GetStringSlice(keys ...string) (slice []string, ok bool) {
	rawSlice, ok := m.GetSlice(keys...)
	if !ok {
		return
	}

	for _, rawElement := range rawSlice {
		if element, ok := rawElement.(string); ok {
			slice = append(slice, element)
		}
	}
	return
}

func (m *DataBag) GetSubOptions(keys ...string) (result DataBag, ok bool) {
	value, ok := tool.GetMapKey(*m, keys...)
	if ok {
		result, ok = value.(map[string]interface{})
	}
	return
}

func (m *DataBag) GetString(keys ...string) (result string, ok bool) {
	value, ok := tool.GetMapKey(*m, keys...)
	if ok {
		result, ok = value.(string)
	}
	return
}

func (m *DataBag) Decode(result interface{}, keys ...string) error {
	value, ok := m.Get(keys...)
	if !ok {
		return fmt.Errorf("key not found %v", strings.Join(keys, "."))
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, result)
}

func convertMapToStringMap(in interface{}) (out interface{}, err error) {
	mapString, ok := in.(map[string]interface{})
	if ok {
		for k, v := range mapString {
			mapString[k], err = convertMapToStringMap(v)
			if err != nil {
				return
			}
		}
		return mapString, nil
	}

	mapInterface, ok := in.(map[interface{}]interface{})
	if ok {
		mapString := make(map[string]interface{})
		for k, v := range mapInterface {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("failed to convert %v to string", k)
			}

			mapString[key], err = convertMapToStringMap(v)
			if err != nil {
				return
			}
		}
		return mapString, nil
	}

	return in, nil
}

func (m *DataBag) Sanitize() (err error) {
	n := make(DataBag)
	for k, v := range *m {
		n[k], err = convertMapToStringMap(v)
		if err != nil {
			return
		}
	}
	*m = n
	return
}

func getOptionsMap(optionKey string, primary, secondary DataBag) (value DataBag) {
	value, ok := primary.GetSubOptions(optionKey)
	if !ok {
		value, _ = secondary.GetSubOptions(optionKey)
	}

	return
}

func getOptions(optionKey string, primary, secondary DataBag) (value []interface{}, ok bool) {
	value, ok = primary.GetSlice(optionKey)
	if !ok {
		value, ok = secondary.GetSlice(optionKey)
	}

	return
}

package setting_test

import (
	"bytes"
	"testing"

	ini "gopkg.in/ini.v1"

	"unicode/utf8"

	"fmt"

	"github.com/pquerna/ffjson/ffjson"

	"strings"

	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/gogits/gogs/pkg/bindata"
)

func TestSettingExport(t *testing.T) {

	Cfg, err := ini.Load(bindata.MustAsset("conf/app.ini"))
	if err != nil {
		t.Fatal(2, "Fail to parse 'conf/app.ini': %v", err)
	}

	Cfg.NameMapper = ini.AllCapsUnderscore

	sections := Cfg.Sections()

	valuesMap := make(map[string]interface{})

	for _, s := range sections {

		var fileName string
		if s.Name() == "DEFAULT" {
			fileName = "global"
		} else {
			fileName = s.Name()
		}

		keys := s.Keys()

		for _, key := range keys {

			if !utf8.ValidString(key.Value()) {
				t.Log("not utf8 code")
			}

			add(fileName+"."+key.Name(), convert(key), valuesMap)

			//valuesMap[key.Name()] = key.Value()
		}

		//
		//fmt.Print(section.ToString())
		//fmt.Print("\n")
	}

	src := &bytes.Buffer{}
	enc := ffjson.NewEncoder(src)
	enc.SetEscapeHTML(false)
	enc.Encode(valuesMap)

	dst := &bytes.Buffer{}
	json.Indent(dst, src.Bytes(), "", "\t")

	fmt.Println(dst.String())

	yml, err := yaml.JSONToYAML(dst.Bytes())

	fmt.Println(string(yml))

}
func add(key string, value interface{}, i map[string]interface{}) {

	keys := strings.Split(key, ".")

	if len(keys) == 1 {
		if _, ok := i[key]; ok {

		} else {
			i[updateKey(key)] = value
		}
	} else {
		AddByPath(keys, i, value)
	}

}

func convert(value *ini.Key) interface{} {
	val, err := value.Int()
	if err == nil {
		return val
	} else {
		val, err := value.Bool()
		if err == nil {
			return val
		} else {
			return value.String()
		}
	}
}

func updateKey(key string) string {

	key = strings.ToLower(key)

	return key
}

func AddByPath(path []string, i map[string]interface{}, value interface{}) {

	if len(path) > 1 {
		if result, ok := i[path[0]]; ok {

			if val, ok := result.(map[string]interface{}); ok {
				AddByPath(path[1:], val, value)
			} else {

				if sval, ok := result.(string); ok {

					tm := make(map[string]interface{})

					key := path[0]

					tm[updateKey(key)+"_name"] = sval

					AddByPath(path[1:], tm, value)

					i[updateKey(key)] = tm
				}

				//fmt.Printf(" val is not map: %s , %s, value %v \n", strings.Join(path, "."), path[0], result)
			}

		} else {
			tm := make(map[string]interface{})
			AddByPath(path[1:], tm, value)

			key := path[0]

			i[updateKey(key)] = tm
		}
	} else {
		if result, ok := i[path[0]]; ok {
			//
			if val, ok := result.(map[string]interface{}); ok {

				key := path[0]

				val[updateKey(key)] = value
			} else {
				fmt.Printf("key: %s , value %v \n", path[0], value)
			}

		} else {

			key := path[0]
			i[updateKey(key)] = value
		}
	}
}

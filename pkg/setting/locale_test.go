package setting

import (
	"bytes"
	"testing"

	ini "gopkg.in/ini.v1"

	"os"

	"unicode/utf8"

	"fmt"

	"github.com/pquerna/ffjson/ffjson"

	"encoding/json"

	"strings"

	"github.com/gogits/gogs/pkg/bindata"
)

func TestLocale(t *testing.T) {

	localFiles := make(map[string][]byte)
	localFiles["en"] = bindata.MustAsset("conf/locale/locale_en-US.ini")
	localFiles["zh-CN"] = bindata.MustAsset("conf/locale/locale_zh-CN.ini")

	for code, v := range localFiles {
		Cfg, err := ini.Load(v)
		if err != nil {
			t.Fatal(2, "Fail to parse 'conf/app.ini': %v", err)
		}

		Cfg.NameMapper = ini.AllCapsUnderscore

		sections := Cfg.Sections()

		for _, s := range sections {

			var fileName string
			if s.Name() == "DEFAULT" {
				fileName = "default"
			} else {
				fileName = s.Name()
			}

			localMap := make(map[string]map[string]interface{})

			valuesMap := make(map[string]interface{})

			keys := s.Keys()
			for _, key := range keys {

				if !utf8.ValidString(key.Value()) {
					t.Log("not utf8 code")
				}

				add(key.Name(), key.Value(), valuesMap)

				//valuesMap[key.Name()] = key.Value()
			}

			localMap[fileName] = valuesMap

			fileToWrite := "/Workspace/Workshop/go/src/github.com/gogits/gogs/i18n"
			if _, err := os.Stat(fileToWrite); os.IsNotExist(err) {
				// path/to/whatever does not exist
				err = os.Mkdir(fileToWrite, 0700)
				if err != nil {
					t.Fatal(2, "Fail to write file  %v", err)
				}
			}
			fileToWrite = fileToWrite + "/" + code
			if _, err := os.Stat(fileToWrite); os.IsNotExist(err) {
				// path/to/whatever does not exist
				err = os.Mkdir(fileToWrite, 0700)
				if err != nil {
					t.Fatal(2, "Fail to write file  %v", err)
				}
			}

			src := &bytes.Buffer{}
			enc := ffjson.NewEncoder(src)
			enc.SetEscapeHTML(false)
			enc.Encode(localMap)

			dst := &bytes.Buffer{}
			json.Indent(dst, src.Bytes(), "", "\t")

			j := dst.Bytes()

			fileToWrite = fileToWrite + "/" + fileName + ".json"
			if _, err := os.Stat(fileToWrite); os.IsNotExist(err) {
				// path/to/whatever does not exist
				f, err := os.Create(fileToWrite)
				if err != nil {
					t.Fatal(2, "Fail to write file  %v", err)
				}

				defer f.Close()

				f.Write(j)

				f.Sync()

			} else {
				f, err := os.OpenFile(fileToWrite, os.O_RDWR, 0600)

				if err != nil {
					t.Fatal(2, "Fail to write file  %v", err)
				}

				defer f.Close()
				f.Write(j)
				//	f.WriteString(html.UnescapeString(string(j)))
				f.Sync()
			}

			//
			//fmt.Print(section.ToString())
			//fmt.Print("\n")
		}
	}

}
func add(key string, value string, i map[string]interface{}) {

	keys := strings.Split(key, ".")

	if len(keys) == 1 {
		if _, ok := i[key]; ok {

		} else {
			i[key] = value
		}
	} else {
		AddByPath(keys, i, value)
	}

}

func AddByPath(path []string, i map[string]interface{}, value string) {

	if len(path) > 1 {
		if result, ok := i[path[0]]; ok {

			if val, ok := result.(map[string]interface{}); ok {
				AddByPath(path[1:], val, value)
			} else {

				if sval, ok := result.(string); ok {

					tm := make(map[string]interface{})

					tm[path[0]+"_name"] = sval

					AddByPath(path[1:], tm, value)

					i[path[0]] = tm
				}

				//fmt.Printf(" val is not map: %s , %s, value %v \n", strings.Join(path, "."), path[0], result)
			}

		} else {
			tm := make(map[string]interface{})
			AddByPath(path[1:], tm, value)
			i[path[0]] = tm
		}
	} else {
		if result, ok := i[path[0]]; ok {
			//
			if val, ok := result.(map[string]interface{}); ok {
				val[path[0]] = value
			} else {
				fmt.Printf("key: %s , value %v \n", path[0], value)
			}

		} else {
			i[path[0]] = value
		}
	}
}

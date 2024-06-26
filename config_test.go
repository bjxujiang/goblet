package goblet

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	toml "github.com/extrame/go-toml-config"
	"github.com/bjxujiang/goblet/config"
	"github.com/extrame/unmarshall"
	myyaml "github.com/extrame/unmarshall/yaml"
	"gopkg.in/yaml.v3"
)

func TestConfigSubStruct(t *testing.T) {
	var obj struct {
		Type   string `goblet:"type,t1"`
		Detail []struct {
			Name string `goblet:"name,here"`
			Sex  int    `goblet:"sex"`
		} `goblet:"array"`
	}

	var node = new(yaml.Node)

	err := yaml.NewDecoder(strings.NewReader(`
basic:
  www_root: ./www3214
test:
  type: t2
  array:
    - name: 1
      sex: 2
    - name: 2
      sex: 3
  `)).Decode(node)

	// for n, m := range node {
	// 	fmt.Println(n, m)
	// }

	if err != nil {
		t.Fatal(err)
	}

	basic, _ := myyaml.GetChildNode(node, "basic")
	var b config.Basic
	basic.Decode(&b)

	fmt.Println(b)

	// ctx := fetch(myyaml.GetChildNode(node, "test"))

	// fmt.Println(ctx)

	// var u = unmarshall.Unmarshaller{
	// 	ValueGetter: func(tag string) []string {
	// 		fmt.Println(tag, ctx[tag])
	// 		if v, ok := ctx[tag]; ok {
	// 			return []string{v}
	// 		} else {
	// 			return []string{}
	// 		}
	// 	},
	// 	ValuesGetter: func(prefix string) url.Values {
	// 		return make(url.Values)
	// 	},
	// 	TagConcatter: func(prefix string, tag string) string {
	// 		return prefix + "." + tag
	// 	},
	// 	Tag:      "goblet",
	// 	AutoFill: true,
	// }
	// u.Unmarshall(&obj)
	fmt.Println(obj)
}

func TestConfigSubArray(t *testing.T) {
	var obj struct {
		Type   string `goblet:"type,t1"`
		Detail []struct {
			Name string `goblet:"name,here"`
		} `goblet:"detail"`
	}

	err := toml.Parse("./test/test_array.config")

	if err != nil {
		t.Fatal(err)
	}

	var u = unmarshall.Unmarshaller{
		ValueGetter: func(tag string) []string {
			pText := toml.String("test"+"."+tag, "")
			toml.Load()
			if *pText != "" {
				return []string{*pText}
			} else {
				return []string{}
			}

		},
		Tag: "goblet",
		ValuesGetter: func(prefix string) url.Values {
			return make(url.Values)
		},
		TagConcatter: func(prefix string, tag string) string {
			return prefix + "." + tag
		},
		BaseName: func(path string, prefix string) string {
			return strings.Split(strings.TrimPrefix(path, prefix+"["), "]")[0]
		},
		AutoFill: true,
	}
	u.Unmarshall(&obj)
	fmt.Println(obj)
}

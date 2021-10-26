package jsontyperegistry_test

import (
	"testing"

	"github.com/theplant/jsontyperegistry"
	"github.com/theplant/testingutils"
)

type Post struct {
	Title  string
	Author *Author
}

type Author struct {
	Name string
	Age  int
}

func TestAll(t *testing.T) {
	for _, c := range cases {
		jsontyperegistry.MustRegisterType(c.value)
	}
	
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			txt := jsontyperegistry.MustJSONString(c.value)
			t.Log(txt)
			v := jsontyperegistry.MustNewWithJSONString(txt)

			diff := testingutils.PrettyJsonDiff(c.value, v)
			if len(diff) > 0 {
				t.Error(diff)
			}
		})
	}

}

var cases = []struct {
	name  string
	value interface{}
}{
	{
		name: "array of objects",
		value: []*Post{
			{
				Author: &Author{
					Name: "Tom",
					Age:  33,
				},
				Title: "123",
			},
			{
				Author: &Author{
					Name: "John",
					Age:  35,
				},
				Title: "456",
			},
		},
	},

	{
		name: "object pointer",
		value: &Post{
			Author: &Author{
				Name: "Tom",
				Age:  33,
			},
			Title: "123",
		},
	},
	{
		name: "object pointer again",
		value: &Post{
			Author: &Author{
				Name: "Tom",
				Age:  33,
			},
			Title: "123",
		},
	},
	{
		name: "object struct",
		value: Post{
			Author: &Author{
				Name: "Tom",
				Age:  33,
			},
			Title: "123",
		},
	},

	{
		name: "map of posts",
		value: map[string]Post{
			"1": {
				Author: &Author{
					Name: "Tom",
					Age:  33,
				},
				Title: "123",
			},
		},
	},

	{
		name:  "simple string array",
		value: []string{"hello"},
	},

	{
		name:  "simple string",
		value: "hello",
	},

	{
		name:  "simple int",
		value: 123,
	},
}

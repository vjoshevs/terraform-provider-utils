package provider

import (
	"reflect"
	"testing"
)

func TestMergeMaps(t *testing.T) {
	cases := []struct {
		dst    map[interface{}]interface{}
		src    map[interface{}]interface{}
		result map[interface{}]interface{}
	}{
		// merge maps
		{
			dst: map[interface{}]interface{}{
				"e1": "abc",
			},
			src: map[interface{}]interface{}{
				"e2": "def",
			},
			result: map[interface{}]interface{}{
				"e1": "abc",
				"e2": "def",
			},
		},
		// merge nested maps
		{
			dst: map[interface{}]interface{}{
				"root": map[interface{}]interface{}{
					"child1": "abc",
				},
			},
			src: map[interface{}]interface{}{
				"root": map[interface{}]interface{}{
					"child2": "def",
				},
			},
			result: map[interface{}]interface{}{
				"root": map[interface{}]interface{}{
					"child1": "abc",
					"child2": "def",
				},
			},
		},
	}

	for _, c := range cases {
		MergeMaps(reflect.ValueOf(c.dst), reflect.ValueOf(c.src), true)
		if !reflect.DeepEqual(c.dst, c.result) {
			t.Fatalf("Error matching dst and result: %#v vs %#v", c.dst, c.result)
		}
	}
}

func TestListItem(t *testing.T) {
	cases := []struct {
		dst            map[interface{}]interface{}
		key            string
		src            interface{}
		mergeListItems bool
		result         map[interface{}]interface{}
	}{
		// merge primitive list items
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
				},
			},
			key:            "list",
			src:            "ghi",
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
					"ghi",
				},
			},
		},
		// merge matching primitive list items
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
				},
			},
			key:            "list",
			src:            "abc",
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
				},
			},
		},
		// merge matching primitive list items
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
				},
			},
			key:            "list",
			src:            "abc",
			mergeListItems: false,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					"abc",
					"def",
				},
			},
		},
		// merge matching map list items
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
						},
					},
				},
			},
			key: "list",
			src: map[interface{}]interface{}{
				"name": "abc",
				"map": map[interface{}]interface{}{
					"elem3": "value3",
				},
			},
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
							"elem3": "value3",
						},
					},
				},
			},
		},
		// append matching map list items
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
						},
					},
				},
			},
			key: "list",
			src: map[interface{}]interface{}{
				"name": "abc",
				"map": map[interface{}]interface{}{
					"elem3": "value3",
				},
			},
			mergeListItems: false,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
						},
					},
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem3": "value3",
						},
					},
				},
			},
		},
		// merge matching map list items with extra src primitive attribute
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name": "abc",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
						},
					},
				},
			},
			key: "list",
			src: map[interface{}]interface{}{
				"name":  "abc",
				"name2": "def",
				"map": map[interface{}]interface{}{
					"elem3": "value3",
				},
			},
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name":  "abc",
						"name2": "def",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
							"elem3": "value3",
						},
					},
				},
			},
		},
		// merge matching map list items with extra dst primitive attribute
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name":  "abc",
						"name2": "def",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
						},
					},
				},
			},
			key: "list",
			src: map[interface{}]interface{}{
				"name": "abc",
				"map": map[interface{}]interface{}{
					"elem3": "value3",
				},
			},
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name":  "abc",
						"name2": "def",
						"map": map[interface{}]interface{}{
							"elem1": "value1",
							"elem2": "value2",
							"elem3": "value3",
						},
					},
				},
			},
		},
		// not merge matching dict list items with extra dst and src primitive attribute
		{
			dst: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name":  "abc",
						"name2": "def",
					},
				},
			},
			key: "list",
			src: map[interface{}]interface{}{
				"name":  "abc",
				"name3": "ghi",
			},
			mergeListItems: true,
			result: map[interface{}]interface{}{
				"list": []interface{}{
					map[interface{}]interface{}{
						"name":  "abc",
						"name2": "def",
					},
					map[interface{}]interface{}{
						"name":  "abc",
						"name3": "ghi",
					},
				},
			},
		},
	}

	for _, c := range cases {
		MergeListItem(reflect.ValueOf(c.dst), reflect.ValueOf(c.key), reflect.ValueOf(c.src), c.mergeListItems)
		if !reflect.DeepEqual(c.dst, c.result) {
			t.Fatalf("Error matching dst and result: %#v vs %#v", c.dst, c.result)
		}
	}
}

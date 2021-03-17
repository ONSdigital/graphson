package graphson

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	genericID = GenericValue{"g:String", "id123"}
	p         = "proppy"
)

func TestStringProperties(t *testing.T) {
	type Expectant struct {
		Label        string
		ExpectLabels []string
		ExpectSingle *string
		ExpectMulti  map[string][]string
		ExpectMeta   map[string]map[string][]string
		ExpectInt32  map[string][]int32
		ExpectInt64  map[string][]int64
		ExpectBool   map[string][]bool
	}
	exStr := map[string][]string{"tSimple": {"tSimple"}, "tMulti": {"tMulti1", "tMulti2"}, "tMeta": {"tMeta1", "tMeta2"}}
	givens := []Vertex{
		// test no properties
		{Type: "g:Vertex",
			Value: VertexValue{ID: "test-id1", Label: "lable1::lable2"},
		},
		// test single property
		{Type: "g:Vertex",
			Value: VertexValue{ID: "test-id1", Label: "lable",
				Properties: map[string][]VertexProperty{
					p: {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: p,
								Value: exStr["tSimple"][0],
							},
						},
					},
					"counter": {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: "counter",
								Value: map[string]interface{}{"@type": "g:Int32", "@value": float64(1234)},
							},
						},
					},
					"big-counter": {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: "big-counter",
								Value: map[string]interface{}{"@type": "g:Int64", "@value": float64(1234)},
							},
						},
					},
					"George": {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: "George",
								Value: true,
							},
						},
					},
				},
			},
		},
		// test multi-property
		{Type: "g:Vertex",
			Value: VertexValue{ID: "test-id1", Label: "lable",
				Properties: map[string][]VertexProperty{
					p: {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: p,
								Value: exStr["tMulti"][0],
							},
						},
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: p,
								Value: exStr["tMulti"][1],
							},
						},
					},
				},
			},
		},
		// test meta-property
		{Type: "g:Vertex",
			Value: VertexValue{ID: "test-id1", Label: "lable",
				Properties: map[string][]VertexProperty{
					p: {
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: p + "1",
								Value: exStr["tMeta"][0],
							},
						},
						{Type: "g:VertexProperty",
							Value: VertexPropertyValue{ID: genericID, Label: p + "2",
								Value: exStr["tMeta"][1],
							},
						},
					},
				},
			},
		},
	}
	expecteds := []Expectant{
		{Label: "tZero",
			ExpectLabels: []string{"lable1", "lable2"},
		},
		{Label: "tSimple",
			ExpectLabels: []string{`lable`},
			ExpectSingle: &exStr["tSimple"][0],
			ExpectMulti:  map[string][]string{p: exStr["tSimple"]},
			ExpectMeta:   map[string]map[string][]string{p: {p: exStr["tSimple"]}},
			ExpectInt32:  map[string][]int32{"counter": {1234}},
			ExpectInt64:  map[string][]int64{"big-counter": {1234}},
			ExpectBool:   map[string][]bool{"George": {true}},
		},
		{Label: "tMulti",
			ExpectLabels: []string{`lable`},
			ExpectSingle: nil,
			ExpectMulti:  map[string][]string{p: exStr["tMulti"]},
			ExpectMeta:   map[string]map[string][]string{p: {p: exStr["tMulti"]}},
		},
		{Label: "tMeta",
			ExpectLabels: []string{`lable`},
			ExpectSingle: nil,
			ExpectMulti:  nil,
			ExpectMeta:   map[string]map[string][]string{p: {p + "1": exStr["tMeta"][0:1], p + "2": exStr["tMeta"][1:2]}},
		},
	}

	for i, given := range givens {
		expected := expecteds[i]

		Convey("Test String Labels: "+expecteds[i].Label, t, func() {
			gotLabels := given.GetLabels()
			So(gotLabels, ShouldResemble, expected.ExpectLabels)
		})

		Convey("Test String Props: "+expecteds[i].Label, t, func() {
			_, err := given.GetMultiProperty("not-there")
			So(err, ShouldEqual, ErrorPropertyNotFound)

			gotVal, err := given.GetProperty(p)
			if expected.ExpectSingle == nil {
				So(err, ShouldNotBeNil)
				if expected.Label == "tMulti" || expected.Label == "tMeta" {
					So(err, ShouldEqual, ErrorPropertyIsMulti)
				} else {
					So(err, ShouldEqual, ErrorPropertyNotFound)
				}
			} else {
				So(err, ShouldBeNil)
				So(gotVal, ShouldEqual, *expected.ExpectSingle)
			}

			for key, expectedVals := range expected.ExpectMulti {
				gotVals, err := given.GetMultiProperty(key)
				So(err, ShouldBeNil)
				So(gotVals, ShouldResemble, expectedVals)
			}

			for key, expectedVals := range expected.ExpectMeta {
				gotMap, err := given.GetMetaProperty(key)
				So(err, ShouldBeNil)
				So(gotMap, ShouldResemble, expectedVals)
			}

			for key, expectedVals := range expected.ExpectInt32 {
				gotVals, err := given.GetMultiPropertyInt32(key)
				So(err, ShouldBeNil)
				So(gotVals, ShouldResemble, expectedVals)
			}

			for key, expectedVals := range expected.ExpectInt64 {
				gotVals, err := given.GetMultiPropertyInt64(key)
				So(err, ShouldBeNil)
				So(gotVals, ShouldResemble, expectedVals)
			}

			for key, expectedVals := range expected.ExpectBool {
				gotVals, err := given.GetMultiPropertyBool(key)
				So(err, ShouldBeNil)
				So(gotVals, ShouldResemble, expectedVals)
			}
		})

	}
}

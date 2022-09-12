package jsonutil_test

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/zainirfan13/graphql-client/internal/jsonutil"
)

func TestUnmarshalGraphQL(t *testing.T) {
	/*
		query {
			me {
				name
				height
			}
		}
	*/
	type query struct {
		Me struct {
			Name   string
			Height float64
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"me": {
			"name": "Luke Skywalker",
			"height": 1.72
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want query
	want.Me.Name = "Luke Skywalker"
	want.Me.Height = 1.72
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_graphqlTag(t *testing.T) {
	type query struct {
		Foo string `graphql:"baz"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"baz": "bar"
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: "bar",
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_jsonTag(t *testing.T) {
	type query struct {
		Foo string `json:"baz"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": "bar"
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: "bar",
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_jsonRawTag(t *testing.T) {
	type query struct {
		Data    json.RawMessage
		Another string
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"Data": { "foo":"bar" },
		"Another" : "stuff"
        }`), &got)

	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Another: "stuff",
		Data:    []byte(`{"foo":"bar"}`),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal: %v %v", want, got)
	}
}

func TestUnmarshalGraphQL_fieldAsScalar(t *testing.T) {
	type query struct {
		Data    json.RawMessage  `scalar:"true"`
		DataPtr *json.RawMessage `scalar:"true"`
		Another string
		Tags    map[string]int `scalar:"true"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
                "Data" : {"ValA":1,"ValB":"foo"},
                "DataPtr" : {"ValC":3,"ValD":false},
		"Another" : "stuff",
                "Tags": {
                    "keyA": 2,
                    "keyB": 3
                }
        }`), &got)

	if err != nil {
		t.Fatal(err)
	}
	dataPtr := json.RawMessage(`{"ValC":3,"ValD":false}`)
	want := query{
		Data:    json.RawMessage(`{"ValA":1,"ValB":"foo"}`),
		DataPtr: &dataPtr,
		Another: "stuff",
		Tags: map[string]int{
			"keyA": 2,
			"keyB": 3,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal: %v %v", want, got)
	}
}

func TestUnmarshalGraphQL_orderedMap(t *testing.T) {
	type query [][2]interface{}
	got := query{
		{"foo", ""},
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": "bar"
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		{"foo", "bar"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal: %v != %v", got, want)
	}
}

func TestUnmarshalGraphQL_orderedMapAlias(t *testing.T) {
	type Update struct {
		Name string `graphql:"name"`
	}
	got := [][2]interface{}{
		{"update0:update(name:$name0)", &Update{}},
		{"update1:update(name:$name1)", &Update{}},
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{
      "update0": {
        "name": "grihabor"
      },
      "update1": {
        "name": "diman"
      }
}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := [][2]interface{}{
		{"update0:update(name:$name0)", &Update{Name: "grihabor"}},
		{"update1:update(name:$name1)", &Update{Name: "diman"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal: %v != %v", got, want)
	}
}

func TestUnmarshalGraphQL_array(t *testing.T) {
	type query struct {
		Foo []string
		Bar []string
		Baz []string
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			"bar",
			"baz"
		],
		"bar": [],
		"baz": null
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []string{"bar", "baz"},
		Bar: []string{},
		Baz: []string(nil),
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

// When unmarshaling into an array, its initial value should be overwritten
// (rather than appended to).
func TestUnmarshalGraphQL_arrayReset(t *testing.T) {
	var got = []string{"initial"}
	err := jsonutil.UnmarshalGraphQL([]byte(`["bar", "baz"]`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"bar", "baz"}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_objectArray(t *testing.T) {
	type query struct {
		Foo []struct {
			Name string
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []struct{ Name string }{
			{"bar"},
			{"baz"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_orderedMapArray(t *testing.T) {
	type query struct {
		Foo [][][2]interface{}
	}
	got := query{
		Foo: [][][2]interface{}{
			{{"name", ""}},
		},
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: [][][2]interface{}{
			{{"name", "bar"}},
			{{"name", "baz"}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_pointer(t *testing.T) {
	s := "will be overwritten"
	foo := "foo"
	type query struct {
		Foo *string
		Bar *string
	}
	var got query
	got.Bar = &s // Test that got.Bar gets set to nil.
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": "foo",
		"bar": null
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: &foo,
		Bar: nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_objectPointerArray(t *testing.T) {
	bar := "bar"
	baz := "baz"
	type query struct {
		Foo []*struct {
			Name *string
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			null,
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []*struct{ Name *string }{
			{&bar},
			nil,
			{&baz},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_orderedMapNullInArray(t *testing.T) {
	type query struct {
		Foo [][][2]interface{}
	}
	got := query{
		Foo: [][][2]interface{}{
			{{"name", ""}},
		},
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			null,
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: [][][2]interface{}{
			{{"name", "bar"}},
			nil,
			{{"name", "baz"}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_pointerWithInlineFragment(t *testing.T) {
	type actor struct {
		User struct {
			DatabaseID uint64
		} `graphql:"... on User"`
		Login string
	}
	type query struct {
		Author actor
		Editor *actor
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"author": {
			"databaseId": 1,
			"login": "test1"
		},
		"editor": {
			"databaseId": 2,
			"login": "test2"
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want query
	want.Author = actor{
		User:  struct{ DatabaseID uint64 }{1},
		Login: "test1",
	}
	want.Editor = &actor{
		User:  struct{ DatabaseID uint64 }{2},
		Login: "test2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_unexportedField(t *testing.T) {
	type query struct {
		foo *string
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{"foo": "bar"}`), new(query))
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "struct field for \"foo\" doesn't exist in any of 1 places to unmarshal"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestUnmarshalGraphQL_multipleValues(t *testing.T) {
	type query struct {
		Foo *string
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{"foo": "bar"}{"foo": "baz"}`), new(query))
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "invalid token '{' after top-level value"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestUnmarshalGraphQL_multipleValuesInOrderedMap(t *testing.T) {
	type query [][2]interface{}
	q := query{{"foo", ""}}
	err := jsonutil.UnmarshalGraphQL([]byte(`{"foo": "bar"}{"foo": "baz"}`), &q)
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "invalid token '{' after top-level value"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestUnmarshalGraphQL_union(t *testing.T) {
	/*
		{
			__typename
			... on ClosedEvent {
				createdAt
				actor {login}
			}
			... on ReopenedEvent {
				createdAt
				actor {login}
			}
		}
	*/
	type actor struct{ Login string }
	type closedEvent struct {
		Actor     actor
		CreatedAt time.Time
	}
	type reopenedEvent struct {
		Actor     actor
		CreatedAt time.Time
	}
	type issueTimelineItem struct {
		Typename      string        `graphql:"__typename"`
		ClosedEvent   closedEvent   `graphql:"... on ClosedEvent"`
		ReopenedEvent reopenedEvent `graphql:"... on ReopenedEvent"`
	}
	var got issueTimelineItem
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"__typename": "ClosedEvent",
		"createdAt": "2017-06-29T04:12:01Z",
		"actor": {
			"login": "shurcooL-test"
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := issueTimelineItem{
		Typename: "ClosedEvent",
		ClosedEvent: closedEvent{
			Actor: actor{
				Login: "shurcooL-test",
			},
			CreatedAt: time.Unix(1498709521, 0).UTC(),
		},
		ReopenedEvent: reopenedEvent{
			Actor: actor{
				Login: "shurcooL-test",
			},
			CreatedAt: time.Unix(1498709521, 0).UTC(),
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_orderedMapUnion(t *testing.T) {
	/*
		{
			__typename
			... on ClosedEvent {
				createdAt
				actor {login}
			}
			... on ReopenedEvent {
				createdAt
				actor {login}
			}
		}
	*/
	actor := [][2]interface{}{{"login", ""}}
	closedEvent := [][2]interface{}{{"actor", actor}, {"createdAt", time.Time{}}}
	reopenedEvent := [][2]interface{}{{"actor", actor}, {"createdAt", time.Time{}}}
	got := [][2]interface{}{
		{"__typename", ""},
		{"... on ClosedEvent", closedEvent},
		{"... on ReopenedEvent", reopenedEvent},
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"__typename": "ClosedEvent",
		"createdAt": "2017-06-29T04:12:01Z",
		"actor": {
			"login": "shurcooL-test"
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := [][2]interface{}{
		{"__typename", "ClosedEvent"},
		{"... on ClosedEvent", [][2]interface{}{
			{"actor", [][2]interface{}{{"login", "shurcooL-test"}}},
			{"createdAt", time.Unix(1498709521, 0).UTC()},
		}},
		{"... on ReopenedEvent", [][2]interface{}{
			{"actor", [][2]interface{}{{"login", "shurcooL-test"}}},
			{"createdAt", time.Unix(1498709521, 0).UTC()},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("not equal:\ngot: %v\nwant: %v", got, want)
		createdAt := got[1][1].([][2]interface{})[1]
		t.Logf("key: %s, type: %v", createdAt[0], reflect.TypeOf(createdAt[1]))
	}
}

// Issue https://github.com/shurcooL/githubv4/issues/18.
func TestUnmarshalGraphQL_arrayInsideInlineFragment(t *testing.T) {
	/*
		query {
			search(type: ISSUE, first: 1, query: "type:pr repo:owner/name") {
				nodes {
					... on PullRequest {
						commits(last: 1) {
							nodes {
								url
							}
						}
					}
				}
			}
		}
	*/
	type query struct {
		Search struct {
			Nodes []struct {
				PullRequest struct {
					Commits struct {
						Nodes []struct {
							URL string `graphql:"url"`
						}
					} `graphql:"commits(last: 1)"`
				} `graphql:"... on PullRequest"`
			}
		} `graphql:"search(type: ISSUE, first: 1, query: \"type:pr repo:owner/name\")"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"search": {
			"nodes": [
				{
					"commits": {
						"nodes": [
							{
								"url": "https://example.org/commit/49e1"
							}
						]
					}
				}
			]
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want query
	want.Search.Nodes = make([]struct {
		PullRequest struct {
			Commits struct {
				Nodes []struct {
					URL string `graphql:"url"`
				}
			} `graphql:"commits(last: 1)"`
		} `graphql:"... on PullRequest"`
	}, 1)
	want.Search.Nodes[0].PullRequest.Commits.Nodes = make([]struct {
		URL string `graphql:"url"`
	}, 1)
	want.Search.Nodes[0].PullRequest.Commits.Nodes[0].URL = "https://example.org/commit/49e1"
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

package tags_test

import (
	ttest "github.com/alextutea/go-table-tests"
	"github.com/alextutea/go-tags"
	"reflect"
	"testing"
)

type UnitTest struct{}

func Test(t *testing.T) {
	ut := UnitTest{}

	t.Run("ParseTag", ut.TestParseTag)
	t.Run("ParseTagString", ut.TestParseTagString)
	t.Run("IsEmpty", ut.TestIsEmpty)
	t.Run("HasKey", ut.TestHasKey)
}

func (ut *UnitTest) TestParseTag(t *testing.T) {
	type TestInput struct {
		TagName string
		Field   reflect.StructField
	}

	type TestStruct struct {
		FieldWithNoTag          interface{}
		FieldWithOtherTags      interface{} `nottesttag:"key1,key2"`
		FieldWithEmptyTag       interface{} `testtag:""`
		FieldWithKeys           interface{} `testtag:"key1,key2"`
		FieldWithOptions        interface{} `testtag:"option1=value1"`
		FieldWithKeysAndOptions interface{} `testtag:"key1,key2,option1=value1"`
	}

	tst := reflect.TypeOf(TestStruct{})
	fieldWithNoTag, _ := tst.FieldByName("FieldWithNoTag")
	fieldWithOtherTags, _ := tst.FieldByName("FieldWithOtherTags")
	fieldWithEmptyTag, _ := tst.FieldByName("FieldWithEmptyTag")
	fieldWithKeys, _ := tst.FieldByName("FieldWithKeys")
	fieldWithOptions, _ := tst.FieldByName("FieldWithOptions")
	fieldWithKeysAndOptions, _ := tst.FieldByName("FieldWithKeysAndOptions")

	testCases := []ttest.Case{
		{
			In: TestInput{TagName: "testtag", Field: fieldWithNoTag},
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: make(map[string]string),
			},
			Desc:             "When passing a field with no tags and a tag",
			ExpectedBehavior: "Should return a tag with empty non-nil Keys and Options fields",
		},
		{
			In: TestInput{TagName: "testtag", Field: fieldWithOtherTags},
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: make(map[string]string),
			},
			Desc:             "When passing a field with some tags and a tag that is not used in the field",
			ExpectedBehavior: "Should return a tag with empty non-nil Keys and Options fields",
		},
		{
			In: TestInput{TagName: "testtag", Field: fieldWithEmptyTag},
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: make(map[string]string),
			},
			Desc:             "When passing a field with a tag that contains an empty string",
			ExpectedBehavior: "Should return a tag with empty non-nil Keys and Options fields",
		},
		{
			In: TestInput{TagName: "testtag", Field: fieldWithKeys},
			ExpectedOut: tags.Tag{
				Keys:    []string{"key1", "key2"},
				Options: make(map[string]string),
			},
			Desc:             "When passing a field with a tag that contains some keys",
			ExpectedBehavior: "Should return a tag with those keys in the Keys slice and an empty non-nil Options field",
		},
		{
			In: TestInput{TagName: "testtag", Field: fieldWithOptions},
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: map[string]string{"option1": "value1"},
			},
			Desc:             "When passing a field with a tag that contains some option",
			ExpectedBehavior: "Should return a tag with that option in the Options map and an empty non-nil Keys field",
		},
		{
			In: TestInput{TagName: "testtag", Field: fieldWithKeysAndOptions},
			ExpectedOut: tags.Tag{
				Keys:    []string{"key1", "key2"},
				Options: map[string]string{"option1": "value1"},
			},
			Desc:             "When passing a field with a tag that contains some keys and some option",
			ExpectedBehavior: "Should return a tag with those keys and options in their corresponding fields",
		},
	}

	for _, tc := range testCases {
		t.Logf("\t%s", tc.Desc)
		in, _ := tc.In.(TestInput)
		out := tags.ParseTag(in.Field, in.TagName)
		ok, msg := tc.Check(out, nil)
		if !ok {
			t.Logf(ttest.FailureMessage(tc.ExpectedBehavior))
			t.Fatal(msg)
		}
		t.Logf(ttest.SuccessMessage(tc.ExpectedBehavior))
	}
}

func (ut *UnitTest) TestParseTagString(t *testing.T) {
	testCases := []ttest.Case{
		{
			In: "key1,key2,option1=value1",
			ExpectedOut: tags.Tag{
				Keys:    []string{"key1", "key2"},
				Options: map[string]string{"option1": "value1"},
			},
			Desc:             "When passing a tag string that contains both keys and options",
			ExpectedBehavior: "Should add them to the correct fields into the newly created tag",
		},
		{
			In: "",
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: make(map[string]string),
			},
			Desc:             "When passing an empty tag string",
			ExpectedBehavior: "Should return a Tag object with a non-nil empty slice for Keys and a non-nil empty map for Options",
		},
		{
			In: "key1,key2",
			ExpectedOut: tags.Tag{
				Keys:    []string{"key1", "key2"},
				Options: make(map[string]string),
			},
			Desc:             "When passing a tag string that contains only keys",
			ExpectedBehavior: "Should add them to the Keys slice of the new Tag and set a non-nil empty map for Options",
		},
		{
			In: "option1=value1",
			ExpectedOut: tags.Tag{
				Keys:    []string{},
				Options: map[string]string{"option1": "value1"},
			},
			Desc:             "When passing a tag string that contains only options",
			ExpectedBehavior: "Should add them to the Options map and set a non-nil empty slice for Keys",
		},
	}

	for _, tc := range testCases {
		t.Logf("\t%s", tc.Desc)
		in, _ := tc.In.(string)
		out := tags.ParseTagStr(in)
		ok, msg := tc.Check(out, nil)
		if !ok {
			t.Logf(ttest.FailureMessage(tc.ExpectedBehavior))
			t.Fatal(msg)
		}
		t.Logf(ttest.SuccessMessage(tc.ExpectedBehavior))
	}
}

func (ut *UnitTest) TestIsEmpty(t *testing.T) {
	testCases := []ttest.Case{
		{
			In:               tags.Tag{},
			ExpectedOut:      true,
			Desc:             "When passing an empty tag with nil values for Keys and Options",
			ExpectedBehavior: "Should return true",
		},
		{
			In: tags.Tag{
				Keys:    []string{},
				Options: make(map[string]string),
			},
			ExpectedOut:      true,
			Desc:             "When passing an empty tag with zero values for Keys and Options",
			ExpectedBehavior: "Should return true",
		},
		{
			In: tags.Tag{
				Keys:    []string{"key1"},
				Options: make(map[string]string),
			},
			ExpectedOut:      false,
			Desc:             "When passing a tag with the zero value for Options and a key",
			ExpectedBehavior: "Should return false",
		},
		{
			In: tags.Tag{
				Keys:    []string{},
				Options: map[string]string{"option1": "value1"},
			},
			ExpectedOut:      false,
			Desc:             "When passing a tag with the zero value for Keys and an option",
			ExpectedBehavior: "Should return false",
		},
		{
			In: tags.Tag{
				Keys:    []string{"key1"},
				Options: map[string]string{"option1": "value1"},
			},
			ExpectedOut:      false,
			Desc:             "When passing a tag with a key and an option",
			ExpectedBehavior: "Should return false",
		},
	}

	for _, tc := range testCases {
		t.Logf("\t%s", tc.Desc)
		in, _ := tc.In.(tags.Tag)
		out := in.IsEmpty()
		ok, msg := tc.Check(out, nil)
		if !ok {
			t.Logf(ttest.FailureMessage(tc.ExpectedBehavior))
			t.Fatal(msg)
		}
		t.Logf(ttest.SuccessMessage(tc.ExpectedBehavior))
	}
}

func (ut *UnitTest) TestHasKey(t *testing.T) {
	type TestInput struct {
		Tag tags.Tag
		Key string
	}

	testCases := []ttest.Case{
		{
			In:               TestInput{Key: "key1", Tag: tags.Tag{Keys: []string{"key1", "key2"}}},
			ExpectedOut:      true,
			Desc:             "When passing a tag with some keys and a key that is contained in the tag Keys slice",
			ExpectedBehavior: "Should return true",
		},
		{
			In:               TestInput{Key: "key3", Tag: tags.Tag{Keys: []string{"key1", "key2"}}},
			ExpectedOut:      false,
			Desc:             "When passing a tag with some keys and a key that is not contained in the tag Keys slice",
			ExpectedBehavior: "Should return false",
		},
		{
			In:               TestInput{Key: "key1", Tag: tags.Tag{}},
			ExpectedOut:      false,
			Desc:             "When passing an empty tag with nil values and a key that is contained in the tag Keys slice",
			ExpectedBehavior: "Should return false",
		},
		{
			In:               TestInput{Key: "key1", Tag: tags.Tag{Keys: []string{}}},
			ExpectedOut:      false,
			Desc:             "When passing an empty tag with a zero value for Keys and a key that is contained in the tag Keys slice",
			ExpectedBehavior: "Should return false",
		},
	}

	for _, tc := range testCases {
		t.Logf("\t%s", tc.Desc)
		in, _ := tc.In.(TestInput)
		out := in.Tag.HasKey(in.Key)
		ok, msg := tc.Check(out, nil)
		if !ok {
			t.Logf(ttest.FailureMessage(tc.ExpectedBehavior))
			t.Fatal(msg)
		}
		t.Logf(ttest.SuccessMessage(tc.ExpectedBehavior))
	}
}

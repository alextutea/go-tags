package tags

import (
	"reflect"
	"strings"
)

type Tag struct {
	Keys    []string
	Options map[string]string
}

func (t Tag) IsEmpty() bool {
	return !t.HasKeys() && !t.HasOptions()
}

func (t Tag) HasKeys() bool {
	return len(t.Keys) > 0
}

func (t Tag) HasOptions() bool {
	return len(t.Options) > 0
}

func (t Tag) HasKey(key string) bool {
	for _, k := range t.Keys {
		if k == key {
			return true
		}
	}
	return false
}

func ParseTag(field reflect.StructField, tag string) Tag {
	return ParseTagStr(field.Tag.Get(tag))
}

func ParseTagStr(str string) Tag {
	t := Tag{
		Keys:    []string{},
		Options: make(map[string]string),
	}

	tokens := strings.Split(str, ",")
	for _, token := range tokens {
		if strings.Contains(token, "=") {
			parts := strings.SplitN(token, "=", 2)
			opt := strings.ToLower(parts[0])
			val := strings.ToLower(parts[1])
			t.Options[opt] = val
			continue
		}
		if token != "" {
			t.Keys = append(t.Keys, token)
		}
	}
	return t
}

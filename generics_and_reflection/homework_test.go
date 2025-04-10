package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(input interface{}) string {
	var builder strings.Builder
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		propTag := field.Tag.Get("properties")
		if propTag == "" {
			continue
		}

		if strings.Contains(propTag, "omitempty") {
			if value.IsZero() {
				continue
			}
			propTag = strings.Split(propTag, ",")[0]
		}

		if builder.Len() > 0 {
			builder.WriteString("\n")
		}

		builder.WriteString(propTag)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprintf("%v", value.Interface()))
	}

	return builder.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

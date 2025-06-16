package main

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"nameOfField"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](obj T) string {
	var lines []string

	objType := reflect.TypeOf(obj)
	objVal := reflect.ValueOf(obj)

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)

		tagString := field.Tag.Get("properties")
		if tagString == "" {
			continue
		}

		tagParts := strings.Split(tagString, ",")
		nameOfField := tagParts[0]
		if nameOfField == "" {
			continue
		}

		value := objVal.Field(i)
		isOmitempty := slices.Contains(tagParts, "omitempty")
		if isOmitempty && value.IsZero() {
			continue
		}

		line := fmt.Sprintf("%s=%v", nameOfField, value.Interface())
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "nameOfField=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "nameOfField=John Doe\nage=30\nmarried=true",
		},
		"test case with isOmitempty structField": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "nameOfField=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for nameOfField, test := range tests {
		t.Run(nameOfField, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}

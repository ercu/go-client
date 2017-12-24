package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// ToValues converts a Content type to a url.Values to use with a Ponzu Go client
func ToValues(p interface{}) (url.Values, error) {
	vals := make(url.Values)

	t := reflect.TypeOf(p)
	v := reflect.Indirect(reflect.ValueOf(p))

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {

				if j == 0 {
					vals.Set(jsonTag, fmt.Sprintf("%v", fieldValue.Index(j)))
				} else {

					vals.Add(jsonTag, fmt.Sprintf("%v", fieldValue.Index(j)))
				}
			}
		default:
			vals.Set(jsonTag, fmt.Sprintf("%v", fieldValue))

		}
	}

	return vals, nil
}

// Target represents required criteria to lookup single content items from the
// Ponzu Content API
type Target struct {
	Type string
	ID   int
}

// ParseReferenceURI is a helper method which accepts a reference path / URI from
// a parent Content type, and retrns a Target containing a content item's Type
// and ID
func ParseReferenceURI(uri string) (Target, error) {
	return parseReferenceURI(uri)
}

func parseReferenceURI(uri string) (Target, error) {
	const prefix = "/api/content?"
	if !strings.HasPrefix(uri, prefix) {
		return Target{}, fmt.Errorf("improperly formatted reference URI: %s", uri)
	}

	uri = strings.TrimPrefix(uri, prefix)

	q, err := url.ParseQuery(uri)
	if err != nil {
		return Target{}, fmt.Errorf("failed to parse reference URI: %s, %v", prefix+uri, err)
	}

	if q.Get("type") == "" {
		return Target{}, fmt.Errorf("reference URI missing 'type' value: %s", prefix+uri)
	}

	if q.Get("id") == "" {
		return Target{}, fmt.Errorf("reference URI missing 'id' value: %s", prefix+uri)
	}

	// convert query id string to int
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		return Target{}, err
	}

	return Target{Type: q.Get("type"), ID: id}, nil
}

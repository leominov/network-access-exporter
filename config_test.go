package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestStringMap_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		in  string
		out StringMap
	}{
		{
			in: "- item",
			out: map[string][]string{
				"all": {
					"item",
				},
			},
		},
		{
			in: `group:
- item`,
			out: map[string][]string{
				"group": {
					"item",
				},
			},
		},
	}
	for _, test := range tests {
		var out StringMap
		err := yaml.Unmarshal([]byte(test.in), &out)
		if !assert.NoError(t, err) {
			continue
		}
		assert.Equal(t, test.out, out)
	}
}

package internal

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGetProfiles(t *testing.T) {
	assert := assert.New(t)
	testConfigPath := filepath.Join("testdata", "config")
	tests := []struct {
		name     string
		path     string
		filter   string
		expected []string
		err      bool
	}{
		{"invalid_path", "./unknown", "", nil, true},
		{"testconfig_path", testConfigPath, "", []string{"default", "test1", "test2"}, false},
		{"no_filter", testConfigPath, "", []string{"default", "test1", "test2"}, false},
		{"specific_profile", testConfigPath, "^test1$", []string{"test1"}, false},
		{"profiles_starting_with_test", testConfigPath, "^test.*", []string{"test1", "test2"}, false},
		{"profiles_ending_in_lt", testConfigPath, ".*lt$", []string{"default"}, false},
		{"profiles_containing_a_digit", testConfigPath, `\d`, []string{"test1", "test2"}, false},
		{"invalid_regex", testConfigPath, "\\", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profiles, err := GetProfiles(tt.path, tt.filter)
			if tt.err {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.ElementsMatch(
				tt.expected,
				profiles,
			)
		})
	}
}

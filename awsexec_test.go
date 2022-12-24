package awsexec

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResult struct {
	profile string
	region  string
}

func TestExec(t *testing.T) {
	tests := []struct {
		name     string
		opt      Options
		fn       ExecFunc
		expected map[string]map[string]testResult
		err      bool
	}{
		{
			name: "test",
			opt: Options{
				Profiles: []string{"one", "two", "three"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			fn:  test,
			err: false,
		},
		{
			name: "test_error",
			opt: Options{
				Profiles: []string{"one", "two", "three", "four"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			expected: nil,
			fn:       testErr,
			err:      true,
		},
	}
	for _, tt := range tests {
		if tt.expected == nil && !tt.err {
			tt.expected = expectedResults(t, tt.opt.Profiles, tt.opt.Regions)
		}
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			results := map[string]map[string]testResult{}
			err := Exec(context.Background(), tt.opt, tt.fn, &results)
			if tt.err {
				assert.ErrorContains(err, "test error")
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, results)
			}
		})
	}
}

func TestExecSlice(t *testing.T) {
	tests := []struct {
		name     string
		opt      Options
		fn       ExecFunc
		expected map[string]map[string][]testResult
		err      bool
	}{
		{
			name: "test",
			opt: Options{
				Profiles: []string{"one", "two", "three"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			fn:  testSlice,
			err: false,
		},
		{
			name: "test_error",
			opt: Options{
				Profiles: []string{"one", "two", "three", "four"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			expected: nil,
			fn:       testErr,
			err:      true,
		},
	}
	for _, tt := range tests {
		if tt.expected == nil && !tt.err {
			tt.expected = expectedResultsSlice(t, tt.opt.Profiles, tt.opt.Regions)
		}
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			results := map[string]map[string][]testResult{}
			err := Exec(context.Background(), tt.opt, tt.fn, &results)
			if tt.err {
				assert.ErrorContains(err, "test error")
			} else {
				assert.NoError(err)
				assert.Equal(tt.expected, results)
			}
		})
	}
}

func test(ctx context.Context, profile, region string) (any, error) {
	return testResult{profile, region}, nil
}

func testSlice(ctx context.Context, profile, region string) (any, error) {
	return []testResult{{profile, region}}, nil
}

func testErr(ctx context.Context, profile, region string) (any, error) {
	return nil, fmt.Errorf("test error")
}

func expectedResults(t *testing.T, profiles, regions []string) map[string]map[string]testResult {
	t.Helper()
	out := map[string]map[string]testResult{}
	for _, profile := range profiles {
		out[profile] = map[string]testResult{}
		for _, region := range regions {
			out[profile][region] = testResult{profile, region}
		}
	}
	return out
}

func expectedResultsSlice(t *testing.T, profiles, regions []string) map[string]map[string][]testResult {
	t.Helper()
	out := map[string]map[string][]testResult{}
	for _, profile := range profiles {
		out[profile] = map[string][]testResult{}
		for _, region := range regions {
			out[profile][region] = []testResult{{profile, region}}
		}
	}
	return out
}

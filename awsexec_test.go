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
		opt      *Options
		fn       ExecFunc
		expected []testResult
		err      bool
	}{
		{
			name: "test",
			opt: &Options{
				Profiles: []string{"one", "two", "three"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			fn:  test,
			err: false,
		},
		{
			name: "test_slice",
			opt: &Options{
				Profiles: []string{"one", "three"},
				Regions:  []string{"d"},
			},
			fn:  testSlice,
			err: false,
		},
		{
			name: "test_error",
			opt: &Options{
				Profiles: []string{"one", "two", "three", "four"},
				Regions:  []string{"a", "b", "c", "d"},
			},
			expected: []testResult{},
			fn:       testErr,
			err:      true,
		},
	}
	for _, tt := range tests {
		if tt.expected == nil {
			tt.expected = expectedResults(t, tt.opt.Profiles, tt.opt.Regions)
		}
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			results := []testResult{}
			err := Exec(context.Background(), tt.opt, tt.fn, &results)
			if tt.err {
				assert.ErrorContains(err, "test error")
			} else {
				assert.NoError(err)
			}
			assert.ElementsMatch(tt.expected, results)
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

func expectedResults(t *testing.T, profiles, regions []string) []testResult {
	t.Helper()
	out := []testResult{}
	for _, profile := range profiles {
		for _, region := range regions {
			out = append(out, testResult{profile, region})
		}
	}
	return out
}

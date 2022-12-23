package awsexec

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

type mockEC2Client struct{}

type testResult struct {
	profile string
	region  string
}

func (m *mockEC2Client) DescribeRegions(ctx context.Context, params *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error) {
	return &ec2.DescribeRegionsOutput{
		Regions: []types.Region{
			{RegionName: aws.String("eu-west-1")},
			{RegionName: aws.String("eu-west-2")},
			{RegionName: aws.String("eu-west-3")},
			{RegionName: aws.String("us-east-1")},
			{RegionName: aws.String("us-east-2")},
			{RegionName: aws.String("ap-northeast-1")},
		},
	}, nil
}

func test(ctx context.Context, profile string, cfg aws.Config) (interface{}, error) {
	return testResult{profile, cfg.Region}, nil
}

func testSlice(ctx context.Context, profile string, cfg aws.Config) (interface{}, error) {
	return []testResult{
		{profile, cfg.Region},
	}, nil
}

func TestExec(t *testing.T) {
	assert := assert.New(t)
	testConfigPath := filepath.Join("internal", "testdata", "config")
	svc := &mockEC2Client{}
	tests := []struct {
		name    string
		options *Options
		len     int
	}{
		{"no_filter", &Options{ConfigPath: testConfigPath}, 18},
		{"profile_filter", &Options{ConfigPath: testConfigPath, ProfileFilter: "^test1$"}, 6},
		{"region_filter", &Options{ConfigPath: testConfigPath, RegionFilter: "eu"}, 9},
		{"profile_region_filter", &Options{ConfigPath: testConfigPath, ProfileFilter: `\d`, RegionFilter: "^us-.*$"}, 4},
		{"filter_without_matches", &Options{ConfigPath: testConfigPath, ProfileFilter: "unkown", RegionFilter: "^us-.*$"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, fn := range []ExecFunc{test, testSlice} {
				results := []testResult{}
				err := Exec(&results, fn, tt.options, svc)
				assert.NoError(err)
				assert.Len(results, tt.len)
				for _, r := range results {
					assert.NotZero(r.profile)
					assert.NotZero(r.region)
				}
			}
		})
	}
}

package internal

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

type mockEC2Client struct{}

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

func TestGetRegions(t *testing.T) {
	assert := assert.New(t)
	svc := &mockEC2Client{}
	tests := []struct {
		name     string
		filter   string
		expected []string
		err      bool
	}{
		{"invalid_regex", "\\", nil, true},
		{"no_filter", "", []string{"eu-west-1", "eu-west-2", "eu-west-3", "us-east-1", "us-east-2", "ap-northeast-1"}, false},
		{"eu", "^eu-.*$", []string{"eu-west-1", "eu-west-2", "eu-west-3"}, false},
		{"number_1", "1", []string{"eu-west-1", "us-east-1", "ap-northeast-1"}, false},
		{"east", "east", []string{"ap-northeast-1", "us-east-1", "us-east-2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			regions, err := GetRegions(tt.filter, svc)
			if tt.err {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.ElementsMatch(tt.expected, regions)
		})
	}
}

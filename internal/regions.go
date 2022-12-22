package internal

import (
	"context"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Client interface {
	DescribeRegions(ctx context.Context, params *ec2.DescribeRegionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeRegionsOutput, error)
}

func GetRegions(filter string, svc EC2Client) ([]string, error) {
	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	out, err := svc.DescribeRegions(ctx, &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false), // exclude disabled regions
	})
	if err != nil {
		return nil, err
	}
	regions := []string{}
	for _, r := range out.Regions {
		region := *r.RegionName
		if re.MatchString(region) {
			regions = append(regions, region)
		}
	}
	return regions, nil
}

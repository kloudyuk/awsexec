package internal

import (
	"context"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func GetRegions(ctx context.Context, profile string, filter string) ([]string, error) {
	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}
	client := ec2.NewFromConfig(cfg)
	in := &ec2.DescribeRegionsInput{
		AllRegions: aws.Bool(false),
	}
	out, err := client.DescribeRegions(ctx, in)
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

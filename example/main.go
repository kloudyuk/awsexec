//go:build exclude

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/kloudyuk/awsexec"
)

// Define the function you want to execute across all profile/region combinations.
// The function signature should match:
// func(ctx context.Context, profile, region string) (any, error)
func getLambdaFunctions(ctx context.Context, profile, region string) (any, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	svc := lambda.NewFromConfig(cfg)
	out, err := svc.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, err
	}
	return out.Functions, nil
}

func main() {
	// Define a variable to collect the results. This should always be
	// map[string]map[string]ANY where ANY is your concrete return type.
	// The map[string]map[string] part is required as the package collates
	// the results in a map of profiles containing a map of regions
	results := map[string]map[string][]types.FunctionConfiguration{}

	// Call the awsexec function providing a context, options, the function to execute and
	// a pointer to your results variable
	ctx := context.Background()
	opts := &awsexec.Options{
		Profiles: []string{"dev", "qa"},
		Regions:  []string{"eu-west-1", "eu-west-2"},
	}
	err := awsexec.Exec(ctx, opts, getLambdaFunctions, &results)

	// Handle any errors
	if err != nil {
		log.Fatal(err)
	}

	// Do whatever with the results
	for profile, regions := range results {
		for region, functions := range regions {
			for _, f := range functions {
				fmt.Println(profile, region, *f.FunctionName)
			}
		}
	}
}

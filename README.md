# awsexec

A golang package to execute a given function against many AWS profile/region combinations concurrently

## Background

When you have a lot of AWS accounts and use multiple regions it can be hard to operate against them efficiently. Imagine you want to find and get info on some AWS resources across 50 accounts and 3 regions? That would involve time consuming manual searching or the creation of scripts to loop through your AWS accounts and regions. Once you multiply the accounts and regions by the number of API calls you might make, these scripts can be slow/inefficient to execute.

This package aims to solve that problem by accepting a function and executing it against any number of AWS profile/region combinations concurrently.

## Features

- fast & efficient using goroutines for concurrent execution
- automatic discovery of AWS profiles & regions
- ability to filter AWS profiles/regions
- high test coverage including testing for race conditions

## Usage

```sh
go get -u github.com/kloudyuk/awsexec
```

### Example

For a complete example of how this package can be used, see the example [here](/example/main.go).

## Options

Options can be provided via the awsexec.Options argument to configure which AWS profiles and regions your function is executed against.

If no options are provided, all profiles from your AWS config (~/.aws/config) will be selected and the enabled regions for that account will be looked up in a call to [ec2.DescribeRegions](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#Client.DescribeRegions).

Regex filters can be provided to filter the profiles and/or regions which are automatically selected e.g.

```go
opts := &awsexec.Options{
  ProfileFilter: `^.*dev$`,
  RegionFilter:  `^eu-west-\d$`,
}
```

Alternatively, you can provide the profiles and/or regions to execute against explicitly e.g.

```go
opts := &awsexec.Options{
  Profiles: []string{"dev", "qa"},
  Regions:  []string{"eu-west-1", "eu-west-2"},
}
```

Note that filters will NOT be applied to explicitly defined profiles/regions.

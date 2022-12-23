package awsexec

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/kloudyuk/awsexec/internal"
)

var errorc = make(chan error)
var errors ExecErrors
var execFunc ExecFunc
var execOptions *Options
var results reflect.Value
var resultsMutex sync.Mutex
var wg sync.WaitGroup

type Options struct {
	ConfigPath    string
	ProfileFilter string
	RegionFilter  string
}

type ExecFunc func(ctx context.Context, profile string, cfg aws.Config) (interface{}, error)

type ExecErrors []error

func (ee ExecErrors) Error() string {
	var err string
	for _, e := range ee {
		err += fmt.Sprintf("%s\n", e)
	}
	return err
}

func Exec(res interface{}, fn ExecFunc, opt *Options, svc ...internal.EC2Client) error {
	results = reflect.ValueOf(res).Elem()
	execFunc = fn
	execOptions = opt
	profiles, err := internal.GetProfiles(opt.ConfigPath, opt.ProfileFilter)
	if err != nil {
		return err
	}
	go func() {
		for {
			err := <-errorc
			errors = append(errors, err)
		}
	}()
	for _, profile := range profiles {
		wg.Add(1)
		go execProfile(profile, svc...)
	}
	wg.Wait()
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func execProfile(profile string, svc ...internal.EC2Client) {
	defer wg.Done()
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(profile))
	if err != nil {
		errorc <- fmt.Errorf("error getting config for profile %s: %w", profile, err)
		return
	}
	cfg.Region = "us-east-1"
	var s internal.EC2Client
	if len(svc) > 0 {
		s = svc[0]
	} else {
		s = ec2.NewFromConfig(cfg)
	}
	regions, err := internal.GetRegions(execOptions.RegionFilter, s)
	if err != nil {
		errorc <- fmt.Errorf("error getting regions for profile %s: %w", profile, err)
		return
	}
	for _, region := range regions {
		wg.Add(1)
		cfg.Region = region
		go execRegion(profile, cfg)
	}
}

func execRegion(profile string, cfg aws.Config) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := execFunc(ctx, profile, cfg)
	if err != nil {
		errorc <- fmt.Errorf("awsexec.ExecFunc error for %s %s: %w", profile, cfg.Region, err)
		return
	}
	if r == nil {
		return
	}
	v := reflect.ValueOf(r)
	resultsMutex.Lock()
	defer resultsMutex.Unlock()
	if v.Kind().String() == "slice" {
		results.Set(reflect.AppendSlice(results, v))
	} else {
		results.Set(reflect.Append(results, v))
	}
}

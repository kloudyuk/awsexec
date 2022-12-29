// awsexec can be used to execute a given function for any
// combination of AWS profiles & regions.
//
// The provided function is executed for each profile/region combination
// concurrently using goroutines.
package awsexec

import (
	"context"
	"sync"

	"github.com/kloudyuk/awsexec/internal"
)

var wg sync.WaitGroup
var execFn ExecFunc

type Options struct {
	ConfigPath    string
	ProfileFilter string
	Profiles      []string
	RegionFilter  string
	Regions       []string
}

// ExecFunc defines the function signature expected for the function passed to Exec()
type ExecFunc func(ctx context.Context, profile, region string) (any, error)

// Exec is the main function a package consumer should call
// 'ctx' is passed through all the go routines and eventually into fn
// 'opt' holds options for selecting profiles & regions
// 'fn' is the function to execute for each profile/region combination
// 'results' is a pointer to an object to collate the results from each execution of fn
func Exec(ctx context.Context, opt Options, fn ExecFunc, results any) error {

	// Store fn in a globally accesible var as it'll never change and this
	// avoids having to pass it between the goroutines
	execFn = fn

	// Initialise result & error structs
	// Use reflection to gain access to the underlying results object
	// This allows collating the results of whatever type the results arg points to
	// so the caller doesn't have to worry about type assertions / convertions
	res := NewResult(results)
	errs := &execErr{sync.Mutex{}, []error{}}

	// If we haven't been given profiles explicitly in opt, get the profiles
	// from the AWS Config file (usually ~/.aws/config)
	profiles := opt.Profiles
	if len(profiles) == 0 {
		p, err := internal.GetProfiles(opt.ConfigPath, opt.ProfileFilter)
		if err != nil {
			return err
		}
		profiles = p
	}

	for _, profile := range profiles {
		wg.Add(1)
		go execProfile(ctx, opt, res, errs, profile)
	}

	wg.Wait()

	if errs.Len() > 0 {
		return errs
	}

	return nil

}

func execProfile(ctx context.Context, opt Options, res *Results, errs *execErr, profile string) {
	defer wg.Done()
	regions := opt.Regions
	if len(regions) == 0 {
		r, err := internal.GetRegions(ctx, profile, opt.RegionFilter)
		if err != nil {
			errs.Add(err)
			return
		}
		regions = r
	}
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			r, err := execFn(ctx, profile, region)
			if err != nil {
				errs.Add(err)
				return
			}
			if r != nil {
				res.Add(profile, region, r)
			}
		}(region)
	}
}

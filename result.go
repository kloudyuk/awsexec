package awsexec

import (
	"reflect"
	"sync"
)

type Results struct {
	sync.Mutex
	t reflect.Type
	v reflect.Value
}

func NewResult(v any) *Results {
	val := reflect.ValueOf(v).Elem()
	t := val.Type().Elem()
	return &Results{sync.Mutex{}, t, val}
}

func (r *Results) Add(profile, region string, v any) {
	r.Lock()
	defer r.Unlock()
	key := reflect.ValueOf(profile)
	val := r.v.MapIndex(key)
	if !val.IsValid() {
		r.v.SetMapIndex(key, reflect.MakeMap(r.t))
		val = r.v.MapIndex(key)
	}
	key = reflect.ValueOf(region)
	val.SetMapIndex(key, reflect.ValueOf(v))
}

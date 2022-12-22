package internal

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

var home string

func init() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}
}

func GetProfiles(path string, filter string) ([]string, error) {
	if path == "" {
		path = filepath.Join(home, ".aws", "config")
	}
	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	awsConfig, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	profiles := []string{}
	for _, s := range awsConfig.SectionStrings() {
		if s == "DEFAULT" {
			continue
		}
		name := strings.TrimSpace(strings.TrimPrefix(s, "profile "))
		if re.MatchString(name) {
			profiles = append(profiles, name)
		}
	}
	return profiles, nil
}

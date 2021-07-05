package config

import (
	"regexp"
	"strings"
)

var (
	AuthPermitConfig AuthPermitAll
	HystrixConfig HystrixCfg
)
type AuthPermitAll struct {
	PermitAll []interface{}
}
type HystrixCfg struct {
	Host string
	Port string
}

func Match(str string) bool {
	if len(AuthPermitConfig.PermitAll) > 0 {
		targetValue := AuthPermitConfig.PermitAll
		for i := 0; i < len(targetValue); i++ {
			s      := targetValue[i].(string)
			res, _ := regexp.MatchString(strings.ReplaceAll(s,"**","(.*?)"),str)
			if res {
				return true
			}
		}
	}
	return false
}
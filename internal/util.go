package internal

import (
	"regexp"
)

func GetSubHost(host string, except string) string {
	reg := regexp.MustCompile("^(.+)." + except)
	matchArr := reg.FindStringSubmatch(host)

	var subHost string
	if len(matchArr) > 0 {
		subHost = matchArr[len(matchArr)-1]
	}
	return subHost
}

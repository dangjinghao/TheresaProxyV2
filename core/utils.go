package core

import "github.com/sirupsen/logrus"

func ComponentLogger(name string) *logrus.Entry {
	return BaseLogger.WithFields(logrus.Fields{"type": "component", "name": name})
}

func InSlice[T comparable](src T, target []T) bool {
	for _, i := range target {
		if i == src {
			return true
		}
	}
	return false
}

func ExistDomain(domain string) bool {

	return ProxySites[domain] != nil || Nicknames[domain] != ""
}

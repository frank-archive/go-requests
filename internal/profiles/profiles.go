// package profiles defines sets of default configurations
// that pipo uses for requesting in http protocol.
// a "profile" means a set of
package profiles

import "github.com/frankli0324/go-requests/internal/client"

type Profile []client.Option

var profiles = map[string]Profile{}

func Get(name string) (Profile, bool) {
	v, ok := profiles[name]
	return v, ok
}

func Register(name string, p Profile) bool {
	if _, ok := profiles[name]; ok {
		return false
	}
	profiles[name] = p
	return true
}

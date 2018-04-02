package hat

import "strings"

func urlJoin(elems ...string) string {
	var u strings.Builder
	for i, e := range elems {
		e = strings.TrimPrefix(e, "/")
		e = strings.TrimSuffix(e, "/")

		u.WriteString(e)
		if i != len(elems)-1 {
			u.WriteString("/")
		}
	}
	return u.String()
}
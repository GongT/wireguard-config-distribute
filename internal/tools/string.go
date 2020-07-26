package tools

import "unicode"

func Ucfirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		return u + str[len(u):]
	}
	return ""
}

func ArrayContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ArrayUnique(s []string) (ret []string) {
	uni := map[string]bool{}

	for _, v := range s {
		if !uni[v] {
			uni[v] = true
			ret = append(ret, v)
		}
	}

	return
}

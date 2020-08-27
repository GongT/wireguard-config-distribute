package tools

import (
	"unicode"
)

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

func ArrayFind(s []string, e string) int {
	for index, a := range s {
		if a == e {
			return index
		}
	}
	return -1
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

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

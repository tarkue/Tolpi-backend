package tolpicheckvalid

import (
	"regexp"
)

var REG_ON_LINKS = `[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)?`

func CheckTolpiValid(text *string) bool {
	matched, _ := regexp.MatchString(REG_ON_LINKS, *text)
	if len(*text) > 1 && len(*text) < 2000 && !matched {
		return true
	}
	return false
}

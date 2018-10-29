package rest

import "strings"

func BuildPath(path []string) string {

	return strings.Join(path, "/")
}

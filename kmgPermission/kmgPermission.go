package kmgPermission

import "strings"

type IsAllower interface {
	IsAllow(args map[string]string) bool
}

// 如果这些判断的某一个失败,则该权限失败,如果都通过,则该权限通过.
type And []IsAllower

func (and And) IsAllow(args map[string]string) bool {
	for _, allower := range and {
		if !allower.IsAllow(args) {
			return false
		}
	}
	return true
}

// 如果这些判断的某一个通过,则该权限通过,如果都不通过,则不通过.
type Or []IsAllower

func (or Or) IsAllow(args map[string]string) bool {
	for _, allower := range or {
		if allower.IsAllow(args) {
			return true
		}
	}
	return false
}

type Not struct {
	IsAllower IsAllower
}

func (not Not) IsAllow(args map[string]string) bool {
	return !not.IsAllower.IsAllow(args)
}

type Prefix string

func (prefix Prefix) IsAllow(args map[string]string) bool {
	return strings.HasPrefix(args["n"], string(prefix))
}

var True IsAllower = tTrue{}

type tTrue struct{}

func (t tTrue) IsAllow(args map[string]string) bool {
	return true
}

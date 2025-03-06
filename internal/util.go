package internal

import (
	"maps"
	"slices"
)

func sortAndDedupe(list []string) []string {
	m := map[string]bool{}
	for _, s := range list {
		m[s] = true
	}

	uniqueKeys := slices.Collect(maps.Keys(m))
	slices.Sort(uniqueKeys)

	return uniqueKeys
}

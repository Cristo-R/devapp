package utils

import (
	"gitlab.shoplazza.site/common/plugin-common/xtypes"
)

// slice1 , slice2 交集
func XIdIntersect(slice1, slice2 []xtypes.XId) []xtypes.XId {
	n := make([]xtypes.XId, 0)
	m := make(map[xtypes.XId]int)

	for _, v := range slice2 {
		m[v]++
	}
	for _, v := range slice1 {
		times, _ := m[v]

		if times >= 1 {
			n = append(n, v)
			m[v]--
		}
	}
	return n
}

// slice1 , slice2 差集
func XIdDifference(slice1, slice2 []xtypes.XId) []xtypes.XId {
	n := make([]xtypes.XId, 0)

	inter := XIdIntersect(slice1, slice2)

	m := make(map[xtypes.XId]int)
	for _, v := range inter {
		m[v]++
	}

	for _, k := range slice1 {
		times, _ := m[k]
		if times == 0 {
			n = append(n, k)
		}
	}
	return n
}

// slice1 , slice2 并集
func XIdUnion(slice1, slice2 []xtypes.XId) []xtypes.XId {
	m := make(map[xtypes.XId]struct{})
	for _, v := range slice1 {
		m[v] = struct{}{}
	}
	for _, v := range slice2 {
		m[v] = struct{}{}
	}

	newSlice := make([]xtypes.XId, 0, len(m))
	for k, _ := range m {
		newSlice = append(newSlice, k)
	}
	return newSlice
}

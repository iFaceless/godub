package utils

type Set struct {
	inner map[interface{}]bool
}

func NewSet(items ...interface{}) *Set {
	set := &Set{
		inner: make(map[interface{}]bool, len(items)),
	}
	for _, item := range items {
		set.inner[item] = true
	}

	return set
}

func (s *Set) Has(item interface{}) bool {
	return s.inner[item]
}

func (s *Set) Len() int {
	return len(s.inner)
}

func (s *Set) StringItems() (keys []string) {
	for k := range s.inner {
		keys = append(keys, k.(string))
	}

	return
}

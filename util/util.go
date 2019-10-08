package util

import (
	"github.com/l3njo/yap/model"
)

// FilterA returns a new slice of articles that satisfy the predicate f.
func FilterA(vs []model.Article, f func(model.Article) bool) []model.Article {
	vsf := []model.Article{}
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterG returns a new slice of galleries that satisfy the predicate f.
func FilterG(vs []model.Gallery, f func(model.Gallery) bool) []model.Gallery {
	vsf := []model.Gallery{}
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterF returns a new slice of flickers that satisfy the predicate f.
func FilterF(vs []model.Flicker, f func(model.Flicker) bool) []model.Flicker {
	vsf := []model.Flicker{}
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

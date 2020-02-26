package util

import (
	"github.com/l3njo/yap/model"
)

// FilterA returns a new slice of articles that satisfy the predicate f.
func FilterA(vs []model.Article, f func(model.Article) bool) []model.Article {
	var vsf []model.Article
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterG returns a new slice of galleries that satisfy the predicate f.
func FilterG(vs []model.Gallery, f func(model.Gallery) bool) []model.Gallery {
	var vsf []model.Gallery
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterF returns a new slice of flickers that satisfy the predicate f.
func FilterF(vs []model.Flicker, f func(model.Flicker) bool) []model.Flicker {
	var vsf []model.Flicker
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// FilterR returns a new slice of reactions that satisfy the predicate f.
func FilterR(vs []model.Reaction, f func(model.Reaction) bool) []model.Reaction {
	var vsf []model.Reaction
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

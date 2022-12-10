package main

import (
	"math/rand"
	"sync"
)

type randomItemGenerator struct {
	titles     []string
	descs      []string
	titleIndex int
	descIndex  int
	mtx        *sync.Mutex
	shuffle    *sync.Once
}

func (r *randomItemGenerator) reset() {
	r.mtx = &sync.Mutex{}
	r.shuffle = &sync.Once{}
	r.titles = []string{
		"go",
		"ts",
		"sh",
		"nx",
	}
	r.descs = []string{
		"GO Projects",
		"TypeScript Projects",
		"Shell Projects",
		"Next JS Projects",
	}
	r.shuffle.Do(func() {
		shuff := func(x []string) {
			rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
		}
		shuff(r.titles)
		shuff(r.descs)
	})
}

func (r *randomItemGenerator) next() item {
	if r.mtx == nil {
		r.reset()
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	i := item{
		title:       r.titles[r.titleIndex],
		description: r.descs[r.descIndex],
	}
	r.titleIndex++
	if r.titleIndex >= len(r.titles) {
		r.titleIndex = 0
	}
	r.descIndex++
	if r.descIndex >= len(r.descs) {
		r.descIndex = 0
	}
	return i
}

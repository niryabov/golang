//go:build !solution

package tparallel

import (
	"sync"
)

type T struct {
	parallelBlock chan bool
	levelBlock    chan bool
	nextFlag      chan bool
	once          *sync.Once
	wg            *sync.WaitGroup
}

func (t *T) Parallel() {
	t.once.Do(func() { t.nextFlag <- true })
	<-t.parallelBlock
}

func (t *T) Run(subtest func(t *T)) {
	sTest := T{t.levelBlock, make(chan bool, 1), make(chan bool, 1), &sync.Once{}, &sync.WaitGroup{}}
	t.wg.Add(1)
	go func() {
		subtest(&sTest)
		sTest.once.Do(func() { sTest.nextFlag <- true })
		close(sTest.levelBlock)
		t.wg.Done()
	}()
	<-sTest.nextFlag
	sTest.wg.Wait()

}

func Run(topTests []func(t *T)) {
	tests := make([]*T, 0)
	wg := sync.WaitGroup{}
	block := make(chan bool, 1)
	for i, f := range topTests {
		tests = append(tests, &T{block, make(chan bool, 1), make(chan bool, 1), &sync.Once{}, &sync.WaitGroup{}})
		t := tests[i]
		wg.Add(1)

		go func() {
			f(t)
			t.once.Do(func() { t.nextFlag <- true })
			wg.Done()
			close(t.levelBlock)
		}()
		<-t.nextFlag
	}

	close(block) //запускаются все параллельные горутины
	wg.Wait()
}

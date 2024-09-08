// minCoins_benchmark_test.go
// go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof или go test -bench=. -cpuprofile=cpu.prof
// go tool pprof cpu.prof или go tool pprof -top cpu.prof
// go tool pprof -top cpu.prof > top10.txt

package main

import (
	"testing"
)

func benchmarkMinCoins(b *testing.B, fn func(int, []int) []int, val int, coins []int) {
	for i := 0; i < b.N; i++ {
		fn(val, coins)
	}
}

func BenchmarkMinCoins(b *testing.B) {
	coins := []int{1, 5, 10, 25, 50, 100, 500, 1000}
	benchmarkMinCoins(b, minCoins, 99999, coins)
}

func BenchmarkMinCoins2(b *testing.B) {
	coins := []int{1, 5, 10, 25, 50, 100, 500, 1000}
	benchmarkMinCoins(b, minCoins2, 99999, coins)
}

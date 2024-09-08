// minCoins_test.go
package main

import (
	"reflect"
	"testing"
)

func TestMinCoins(t *testing.T) {
	tests := []struct {
		val    int
		coins  []int
		result []int
	}{
		{13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
		{12, []int{1, 3, 4, 7, 13, 15}, []int{7, 4, 1}},
	}

	for _, test := range tests {
		res := minCoins(test.val, test.coins)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("minCoins(%d, %v) = %v; want %v", test.val, test.coins, res, test.result)
		}
	}
}

func TestMinCoins2(t *testing.T) {
	tests := []struct {
		val    int
		coins  []int
		result []int
	}{
		{13, []int{1, 5, 10}, []int{10, 1, 1, 1}},
		{12, []int{1, 3, 4, 7, 13, 15}, []int{7, 4, 1}},
		{0, []int{1, 2, 5}, []int{}},
		{11, []int{1, 5, 5, 10}, []int{10, 1}},
		{9, []int{5, 1, 3}, []int{5, 3, 1}},
	}

	for _, test := range tests {
		res := minCoins2(test.val, test.coins)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("minCoins2(%d, %v) = %v; want %v", test.val, test.coins, res, test.result)
		}
	}
}

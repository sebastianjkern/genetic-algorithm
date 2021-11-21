package main

import "math/rand"

func RandomBool(likeliness float64) bool {
	if (likeliness < 0) || (likeliness > 100) {
		return false
	}

	f := rand.Float64()
	return likeliness >= f
}

func RandomIntBtw(min int, max int) int {
	return rand.Intn(max-min) + min
}

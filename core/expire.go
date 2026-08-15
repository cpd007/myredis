package core

import (
	"log"
	"time"
)

func DeleteExpiredKeys() {

	for {
		frac := deleteSample()
		if frac < 0.25 {
			break
		}
	}

	log.Println("Deleted keys about to expire. Total Keys: ", len(store))
}

func deleteSample() float64 {

	limit := 20
	deletedKeys := 0

	// assuming that iteration of golang's hash table is randomised
	for k, v := range store {

		if v.ExpiresAt > 0 && v.ExpiresAt < time.Now().UnixMilli() {
			Del(k)
			deletedKeys++
		}

		limit--
		if limit == 0 {
			break
		}
	}

	return float64(deletedKeys) / float64(limit)
}

package util

import (
	"math/rand"

	"github.com/42LoCo42/emo/api"
)

func RandomStat(stats *[]api.Stat) *api.Stat {
	var sum int64
	for _, stat := range *stats {
		sum += stat.Count + stat.Boost
	}

	choice := rand.Int63n(sum)

	for _, stat := range *stats {
		choice -= stat.Count + stat.Boost
		if choice < 0 {
			return &stat
		}
	}

	panic("unreachable")
}

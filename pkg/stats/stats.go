package stats

import (
	"fmt"
	"math"
	"time"
)

/**
Prints the processing rate of a channel. Ticks increments the count by at most 255 per tick
 */
func PrintRate(finished <-chan bool, increment <-chan bool) {
	start := time.Now()
	ticks := 0
	//total := float64(numWords)

	for {
		select {
		case <-increment:
			ticks += 1

			if ticks%500 == 0 {
				curr := float64(ticks)
				//ratio := math.Min(math.Max(curr/total, 0), 1)
				//percent := int32(math.Floor(ratio * 100))
				elapsed := float64(time.Since(start).Seconds())
				rate := curr / elapsed

				//fmt.Printf("\r%d%% %d/wps", percent, int32(rate))
				fmt.Printf("\r%d/wps", int32(rate))
			}

		case <-finished:
			return
		}
	}
}

/**
Prints the processing rate of a channel. Ticks increments the count by at most 255 per tick
 */
func PrintProgress(finished <-chan bool, increment <-chan bool, maxTicks int) {
	start := time.Now()
	ticks := 0
	total := float64(maxTicks)

	for {
		select {
		case <-increment:
			ticks += 1

			if ticks%500 == 0 {
				curr := float64(ticks)
				ratio := math.Min(math.Max(curr/total, 0), 1)
				percent := int32(math.Floor(ratio * 100))
				elapsed := float64(time.Since(start).Seconds())
				rate := curr / elapsed

				fmt.Printf("\r%d%% %d/wps", percent, int32(rate))
			}

		case <-finished:
			return
		}
	}
}

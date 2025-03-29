package calma_test

import (
	"testing"
	"time"

	"github.com/ddddddO/calma"
)

func BenchmarkCalendar_Concurrent_String(b *testing.B) {
	b.ResetTimer()

	now := time.Now()
	for i := 0; i < b.N; i++ {
		c, err := calma.NewCalendarConcurrency(now)
		if err != nil {
			b.Fatal(err)
		}
		cc := c.String()
		_ = cc
	}
}

func BenchmarkCalendar_NoConcurrent_String(b *testing.B) {
	b.ResetTimer()

	now := time.Now()
	for i := 0; i < b.N; i++ {
		c, err := calma.NewCalendar(now)
		if err != nil {
			b.Fatal(err)
		}
		cc := c.String()
		_ = cc
	}
}

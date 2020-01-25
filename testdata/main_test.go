package testdata

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain1(t *testing.T) {
	fmt.Fprintf(os.Stdout, "STDOUT\n")
	fmt.Fprintf(os.Stderr, "STDERR\n")
}

func TestLog(t *testing.T) {
	t.Log("log something")
}

func TestMain3(t *testing.T) {
	t.Parallel()
	time.Sleep(time.Second)
}

func BenchmarkPass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		time.Sleep(time.Millisecond)
	}
}

func BenchmarkAllocs(b *testing.B) {
	b.ReportAllocs()
}

func BenchmarkError(b *testing.B) {
	b.ReportAllocs()
	b.Error("something went wrong")
}

func BenchmarkSkip(b *testing.B) {
	b.ReportAllocs()
	b.Skip("reason why")
}

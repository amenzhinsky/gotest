package tests

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

func BenchmarkFib(b *testing.B) {
	b.ReportAllocs()
	z := 0
	for i := 0; i < b.N; i++ {
		z += i
	}
}

func BenchmarkFib2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		time.Sleep(10 * time.Millisecond)
	}
	b.Error("asdf")
}

func BenchmarkFib3(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		time.Sleep(10 * time.Millisecond)
	}
}

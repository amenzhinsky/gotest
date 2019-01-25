package tests

import (
	"fmt"
	"os"
	"testing"
)

func TestMain1(t *testing.T) {
	fmt.Fprintf(os.Stdout, "STDOUT\n")
	fmt.Fprintf(os.Stderr, "STDERR\n")
}

func TestLog(t *testing.T) {
	t.Log("log something")
}

func TestMain2(t *testing.T) {
	//t.Fatal("failed")
}

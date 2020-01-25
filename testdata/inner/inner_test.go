package inner

import "testing"

func TestSuccess(t *testing.T) {

}

func TestFailure(t *testing.T) {
	t.Fatal("something went wrong")
}

func TestSkip(t *testing.T) {
	t.Skip("let's skip it for now")
}

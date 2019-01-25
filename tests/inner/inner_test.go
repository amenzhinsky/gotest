package inner

import "testing"

func TestFoo(t *testing.T) {
	t.Parallel()
	t.Skip()
}

func TestFail(t *testing.T) {
	t.Fatal("foo")
}

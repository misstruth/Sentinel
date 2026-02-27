package testutil

// AssertEqual 断言相等
func AssertEqual(t interface {
	Helper()
	Errorf(string, ...interface{})
}, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

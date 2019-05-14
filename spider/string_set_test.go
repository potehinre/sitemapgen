package spider

import "testing"

func TestStringSet(t *testing.T) {
	visitedUrls := NewStringSet()
	testUrl := "http://example.com"
	anotherUrl := "http://anotherexample.com"
	visitedUrls.Add(testUrl)
	if !visitedUrls.IsExists(testUrl) {
		t.Fail()
	}
	if visitedUrls.IsExists(anotherUrl) {
		t.Fail()
	}
	all := visitedUrls.All()
	if !(len(all) == 1 && all[0] == testUrl) {
		t.Fail()
	}
}

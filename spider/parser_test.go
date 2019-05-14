package spider

import (
	"testing"
)

func testSliceEq(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestParseFromATags(t *testing.T) {
	text := `
	<body>
	  <a href="one">one</a>
	  <a href="two">two</a>
	</body>
	`
	urls := parseUrls([]byte(text))
	expectedUrls := []string{"one", "two"}
	if !testSliceEq(urls, expectedUrls) {
		t.Fail()
	}
}

func TestParseFromATagsNested(t *testing.T) {
	text := `
	<body>
	  <a href="one">one</a>
	  <a href="two">two</a>
	  <div class="f">
	  	<a rel="b" href="three">three</a>
		<a rel="c" href="four">four</a>
	  </div>
	</body>
	`
	urls := parseUrls([]byte(text))
	expectedUrls := []string{"one", "two", "three", "four"}
	if !testSliceEq(urls, expectedUrls) {
		t.Fail()
	}
}

func TestParseFromATagsBrokenHTML(t *testing.T) {
	text := `
	<body>
	  <a href="one">one</a></a>
	  <a href="two">two</a></a>
	  	<a rel="b" href="three">three</a>
		<a rel="c" href="four">four</a>
	  </div>
	`
	urls := parseUrls([]byte(text))
	expectedUrls := []string{"one", "two", "three", "four"}
	if !testSliceEq(urls, expectedUrls) {
		t.Fail()
	}
}

func TestParseAWithBase(t *testing.T) {
	text := `
	<head>
	    <base href="base/" target="_blank">
	</head>
	<body>
	  <a href="one">one</a>
	  <a href="two">two</a>
	</body>
	`
	urls := parseUrls([]byte(text))
	expectedUrls := []string{"base/one", "base/two"}
	if !testSliceEq(urls, expectedUrls) {
		t.Fail()
	}
}

func TestParseAnchors(t *testing.T) {
	text := `
	<body>
	  <a href="#one">one</a>
	  <a href="#two">two</a>
	</body>
	`
	urls := parseUrls([]byte(text))
	expectedUrls := []string{"#one", "#two"}
	if !testSliceEq(urls, expectedUrls) {
		t.Fail()
	}
}

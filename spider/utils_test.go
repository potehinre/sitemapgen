package spider

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateUrlDomainDiffers(t *testing.T) {
	url := "http://example.com/var"
	vu := NewStringSet()
	if validateUrl(url, "example2.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlAlreadyVisited(t *testing.T) {
	url := "http://example.com/var"
	vu := NewStringSet()
	vu.Add(url)
	if validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlWithFragment(t *testing.T) {
	url := "http://example.com/foo/#comments"
	vu := NewStringSet()
	if validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlIncorrectExt(t *testing.T) {
	url := "http://example.com/foo/bar.jpeg"
	vu := NewStringSet()
	if validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlAnchor(t *testing.T) {
	url := "#bar"
	vu := NewStringSet()
	if validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlValid(t *testing.T) {
	url := "http://example.com/var"
	vu := NewStringSet()
	if !validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestValidateUrlValidHtmlExt(t *testing.T) {
	url := "http://example.com/var.html"
	vu := NewStringSet()
	if !validateUrl(url, "example.com", vu) {
		t.Fail()
	}
}

func TestGetPageOK(t *testing.T) {
	path := "/some/path"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() == path {
			rw.Write([]byte(`OK`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()
	client := server.Client()
	body, err := getPage(client, server.URL+path)
	if err != nil {
		t.Fail()
	}
	if !(string(body) == "OK") {
		t.Fail()
	}
}

func TestGetPage500(t *testing.T) {
	path := "/some/path"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		if req.URL.String() == path {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`ERRR OCCURED`))
		}
	}))
	// Close the server when test finishes
	defer server.Close()
	client := server.Client()
	_, err := getPage(client, server.URL+path)
	if !(err == ErrNotOKStatus) {
		t.Fail()
	}
}

func TestNormalizeUrlStartsWithSlash(t *testing.T) {
	url := "/some/url"
	baseUrl := "http://www.example.com"
	nUrl := normalizeUrl(url, baseUrl)
	if nUrl != "http://www.example.com/some/url" {
		t.Fail()
	}
}

func TestNormalizeUrlCutTrailingSlashes(t *testing.T) {
	url := "http://www.example.com/some/url/"
	baseUrl := "http://www.example.com"
	nUrl := normalizeUrl(url, baseUrl)
	if nUrl != "http://www.example.com/some/url" {
		t.Fail()
	}
}

func TestNormalizeUrlCutTrailingSlashesWithBaseUrl(t *testing.T) {
	url := "/some/url/"
	baseUrl := "http://www.example.com"
	nUrl := normalizeUrl(url, baseUrl)
	if nUrl != "http://www.example.com/some/url" {
		t.Fail()
	}
}

func TestBaseUrlWithoutPort(t *testing.T) {
	url := "http://example.com/some/url/"
	base, err := baseUrl(url)
	if err != nil {
		t.Fail()
	}
	if base != "http://example.com" {
		t.Fail()
	}
}

func TestBaseUrlWithPort(t *testing.T) {
	url := "http://example.com:8000/some/url/"
	base, err := baseUrl(url)
	if err != nil {
		t.Fail()
	}
	if base != "http://example.com:8000" {
		t.Fail()
	}
}

func TestBaseUrlHttpsWithPort(t *testing.T) {
	url := "https://example.com:8000/some/url/"
	base, err := baseUrl(url)
	if err != nil {
		t.Fail()
	}
	if base != "https://example.com:8000" {
		t.Fail()
	}
}

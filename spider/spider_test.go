package spider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Response struct {
	Status int
	Refs   []string
}

func makeWebSite(site map[string]*Response) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		reqUrl := req.URL.String()
		if resp, ok := site[reqUrl]; !ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte(`NotFound`))
		} else {
			rw.WriteHeader(resp.Status)
			hrefs := []string{}
			for _, ref := range resp.Refs {
				hrefs = append(hrefs, fmt.Sprintf(`<a href="%s">Ref</a>`, ref))
			}
			rw.Write([]byte(strings.Join(hrefs, "")))
		}
	}))
	// Close the server when test finishes
	client := server.Client()
	return server, client
}

func isSetsEq(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	firstS := map[string]bool{}
	secondS := map[string]bool{}
	for _, v := range first {
		firstS[v] = true
	}
	for _, v := range second {
		secondS[v] = true
	}
	for k, _ := range firstS {
		if _, ok := secondS[k]; !ok {
			return false
		}
	}
	return true
}

func urlsWithBase(baseUrl string, urls []string) []string {
	res := []string{}
	for _, url := range urls {
		res = append(res, baseUrl+url)
	}
	return res
}

func TestSpiderDepthOne(t *testing.T) {
	webSite := map[string]*Response{
		"/":       {200, []string{"/second", "/third"}},
		"/second": {200, []string{"/fourth", "/fifth"}},
		"/third":  {200, []string{"/seventh", "/eight"}},
	}
	server, client := makeWebSite(webSite)
	defer server.Close()
	urls := Run(client, server.URL, 1, 8)
	if !isSetsEq(urls, urlsWithBase(server.URL, []string{"", "/second", "/third"})) {
		t.Fail()
	}
}

func TestSpiderDepthTwo(t *testing.T) {
	webSite := map[string]*Response{
		"/":       {200, []string{"/second", "/third"}},
		"/second": {200, []string{"/fourth"}},
		"/third":  {200, []string{"/fifth"}},
		"/fourth": {200, []string{}},
		"/fifth":  {200, []string{}},
	}
	server, client := makeWebSite(webSite)
	defer server.Close()
	urls := Run(client, server.URL, 2, 8)
	if !isSetsEq(urls, urlsWithBase(server.URL, []string{"", "/second", "/third", "/fourth", "/fifth"})) {
		t.Fail()
	}
}

func TestSpiderDepthTwoIncorrectStatusCodes(t *testing.T) {
	webSite := map[string]*Response{
		"/":       {200, []string{"/second", "/third"}},
		"/second": {200, []string{"/fourth"}},
		"/third":  {500, []string{"/fifth"}},
		"/fourth": {200, []string{}},
		"/fifth":  {500, []string{}},
	}
	server, client := makeWebSite(webSite)
	defer server.Close()
	urls := Run(client, server.URL, 2, 8)
	if !isSetsEq(urls, urlsWithBase(server.URL, []string{"", "/second", "/fourth"})) {
		t.Fail()
	}
}

func TestSpiderAvoidsRecursion(t *testing.T) {
	webSite := map[string]*Response{
		"/":       {200, []string{"/second", "/third"}},
		"/second": {200, []string{"/"}},
		"/third":  {200, []string{"/fifth"}},
		"/fifth":  {200, []string{}},
	}
	server, client := makeWebSite(webSite)
	defer server.Close()
	urls := Run(client, server.URL, 10, 8)
	if !isSetsEq(urls, urlsWithBase(server.URL, []string{"", "/second", "/third", "/fifth"})) {
		t.Fail()
	}
}

func TestSpiderAvoidsRecursionSamePage(t *testing.T) {
	webSite := map[string]*Response{
		"/":       {200, []string{"/second", "/third"}},
		"/second": {200, []string{"/second"}},
		"/third":  {200, []string{"/fifth"}},
		"/fifth":  {200, []string{}},
	}
	server, client := makeWebSite(webSite)
	defer server.Close()
	urls := Run(client, server.URL, 10, 8)
	if !isSetsEq(urls, urlsWithBase(server.URL, []string{"", "/second", "/third", "/fifth"})) {
		t.Fail()
	}
}

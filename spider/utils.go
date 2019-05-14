package spider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const slash = "/"
const anchor = "#"

const htmlExt = "html"

var ErrNotOKStatus = errors.New("resp status is not ok")

func baseUrl(urlVal string) (string, error) {
	u, err := url.Parse(urlVal)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), nil
}

func normalizeUrl(url string, baseUrl string) string {
	resultUrl := url
	if strings.HasPrefix(url, slash) {
		resultUrl = baseUrl + resultUrl
	}
	return strings.TrimRight(resultUrl, slash)
}

func validateUrl(urlToVal string, startDomain string, visited *StringSet) bool {
	if strings.HasPrefix(urlToVal, anchor) {
		return false
	}
	u, err := url.Parse(urlToVal)
	if err != nil {
		return false
	}
	domain := u.Host
	if domain != startDomain {
		return false
	}
	if visited.IsExists(urlToVal) {
		return false
	}
	if len(u.Fragment) > 0 {
		return false
	}
	sp := strings.Split(u.Path, ".")
	if len(sp) > 1 {
		ext := sp[len(sp)-1]
		if ext != htmlExt {
			return false
		}
	}
	return true
}

func getPage(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error getting url %s:%s", url, err)
		return []byte(""), err
	}
	if resp.StatusCode >= 300 {
		return []byte(""), ErrNotOKStatus
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading url response body %s:%s", url, err)
		return []byte(""), err
	}
	return body, nil
}

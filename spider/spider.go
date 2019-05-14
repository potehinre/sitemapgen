package spider

import (
	"log"
	"net/http"
	"net/url"
)

const workChanCap = 100000

type SpiderTask struct {
	Url   string
	Depth int
}

func Run(client *http.Client, startUrl string, maxdepth int, workerCount int) []string {
	visitedUrls := NewStringSet()
	u, err := url.Parse(startUrl)
	if err != nil {
		log.Fatalf("Error parsing start url %s: %s", startUrl, err)
	}
	startDomain := u.Host
	startBaseUrl, _ := baseUrl(startUrl)

	taskCh := make(chan *SpiderTask, workChanCap)
	stopCh := make(chan bool)
	nStartUrl := normalizeUrl(startUrl, startBaseUrl)
	tt := NewTaskTracker(nStartUrl, stopCh, taskCh)
	for i := 0; i < workerCount; i++ {
		go func() {
			for task := range taskCh {
				func(task *SpiderTask) {
					tt.Running(task.Url)
					urlsToDo := map[string]bool{}
					defer tt.Done(task.Url)
					if visitedUrls.IsExists(task.Url) {
						tt.Done(task.Url)
						return
					}
					log.Printf("Getting url %s width depth %d\n", task.Url, task.Depth)
					body, err := getPage(client, task.Url)
					if err != nil {
						return
					}
					visitedUrls.Add(task.Url)
					if task.Depth >= maxdepth {
						return
					}
					urls := parseUrls(body)
					for _, url := range urls {
						nUrl := normalizeUrl(url, startBaseUrl)
						if validateUrl(nUrl, startDomain, visitedUrls) {
							urlsToDo[nUrl] = true
						}
					}
					tt.DoneWithNewTasks(task.Url, urlsToDo, task.Depth)
				}(task)
			}
		}()
	}
	<-stopCh
	close(taskCh)
	return visitedUrls.All()
}

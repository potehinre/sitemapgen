package spider

import "testing"

func TestTaskTrackerStopEverythingDone(t *testing.T) {
	stopCh := make(chan bool, 10)
	taskCh := make(chan *SpiderTask, 10)
	url := "http://example.com"
	tt := NewTaskTracker(url, stopCh, taskCh)
	tt.Running(url)
	tt.Done(url)
	if !(len(stopCh) == 1) {
		t.Fail()
	}
}

func TestTaskTrackerStopUsualRoutine(t *testing.T) {
	stopCh := make(chan bool, 10)
	taskCh := make(chan *SpiderTask, 100)
	url := "http://example.com"
	tt := NewTaskTracker(url, stopCh, taskCh)
	url2 := "http://example.com/foo"
	url3 := "http://example.com/bar"
	tt.Running(url)
	tt.DoneWithNewTasks(url, map[string]bool{url2: true, url3: true}, 1)
	tt.Running(url2)
	tt.Done(url2)
	tt.Running(url3)
	tt.Done(url3)
	if !(len(stopCh) == 1) {
		t.Fail()
	}
}

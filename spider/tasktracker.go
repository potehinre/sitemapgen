package spider

import "sync"

type TaskStatus int

const (
	ToDo    TaskStatus = iota
	Running TaskStatus = iota
)

func NewTaskTracker(startTask string, stopCh chan<- bool, taskCh chan<- *SpiderTask) *TaskTracker {
	m := map[string]TaskStatus{startTask: ToDo}
	tt := &TaskTracker{status: m, stopCh: stopCh, taskCh: taskCh}
	tt.taskCh <- &SpiderTask{startTask, 0}
	return tt
}

type TaskTracker struct {
	status map[string]TaskStatus
	stopCh chan<- bool
	taskCh chan<- *SpiderTask
	mu     sync.Mutex
}

func (ts *TaskTracker) Running(url string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.status[url] = Running
}

func (ts *TaskTracker) Done(url string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.status, url)
	ts.StopIfEmpty()
}

func (ts *TaskTracker) DoneWithNewTasks(url string, newUrls map[string]bool, prevDepth int) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.status, url)
	for url, _ := range newUrls {
		ts.status[url] = ToDo
		ts.taskCh <- &SpiderTask{url, prevDepth + 1}
	}
	ts.StopIfEmpty()
}

func (ts *TaskTracker) StopIfEmpty() {
	if len(ts.status) == 0 {
		ts.stopCh <- true
	}
}

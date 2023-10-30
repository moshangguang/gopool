package gopool

import "testing"

func TestExample(t *testing.T) {
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	rspChan := make(chan string, len(requests))
	pool := New()
	for _, r := range requests {
		r := r
		pool.Submit(func() {
			rspChan <- r
		})
	}
	pool.StopWait()
	close(rspChan)
	rspSet := map[string]struct{}{}
	for rsp := range rspChan {
		rspSet[rsp] = struct{}{}
	}
	if len(rspSet) < len(requests) {
		t.Fatal("Did not handle all requests")
	}
	for _, req := range requests {
		if _, ok := rspSet[req]; !ok {
			t.Fatal("Missing expected values:", req)
		}
	}
}

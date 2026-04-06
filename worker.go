package main

import "sync"

const WORKER_CNT int = 5
const JQUEUE_LEN int = 20

type Job func()

func add_worker(id int, job_queue <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range job_queue {
		j()
	}

}

func init_job_pool(cnt int, f func()) {
	jqueue := make(chan Job, JQUEUE_LEN)
	var wg sync.WaitGroup
	for i := 0; i < WORKER_CNT; i++ {
		wg.Add(1)
		add_worker(i, jqueue, &wg)
	}

	for i := 0; i < cnt; i++ {
		jqueue <- f
	}
	close(jqueue)
	wg.Wait()
}

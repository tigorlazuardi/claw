package claw

import "sync"

type tracker struct {
	sync.Map
}

func (qt *tracker) Add(jobID int64) {
	qt.Store(jobID, struct{}{})
}

func (qt *tracker) Remove(jobID int64) {
	qt.Delete(jobID)
}

func (qt *tracker) Exists(jobID int64) bool {
	_, exists := qt.Load(jobID)
	return exists
}

func (qt *tracker) List() []int64 {
	var jobs []int64
	qt.Range(func(key, value any) bool {
		if id, ok := key.(int64); ok {
			jobs = append(jobs, id)
		}
		return true
	})
	return jobs
}

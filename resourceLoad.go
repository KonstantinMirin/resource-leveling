package main

type ResourceLoad struct {
	resource    *Resource
	loadSummary map[int]float64
	load        map[int]map[int]float64 //[intervalNo][taskId]hoursWorked
}

func (rl *ResourceLoad) free(interval int) float64 {
	return rl.resource.Capacity - rl.loadSummary[interval]
}

func (rl *ResourceLoad) work(interval, taskId int, time float64) {
	rl.loadSummary[interval] += time
	_, ok := rl.load[interval]
	if !ok {
		rl.load[interval] = make(map[int]float64)
	}
	rl.load[interval][taskId] = time
}

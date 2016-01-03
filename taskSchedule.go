package main
import (
	"fmt"
	"strings"
	"sort"
)

type TaskSchedule struct {
	task         *Task
	finished     bool
	workRemained float64
	workload     map[int]map[int]float64 //[intervalNo][resourceId]hoursWorked
}

func (ts *TaskSchedule) prerequisitesFinished(s *Schedule) bool {
	if ts.task.prerequisites == nil {
		return true
	}
	finishedAll := true
	for i := 0; i < len(ts.task.prerequisites); i++ {
		finishedAll = finishedAll && s.tasks[ts.task.prerequisites[i].Id].finished
	}
	return finishedAll
}

func (ts *TaskSchedule) work(interval, resourceId int, time float64) {
	ts.workRemained -= time
	if ts.workRemained == 0 {
		ts.finished = true
	}
	_, ok := ts.workload[interval]
	if !ok {
		ts.workload[interval] = make(map[int]float64)
	}
	ts.workload[interval][resourceId] = time
}

func (ts TaskSchedule) String() string {
	foundItems := 0
	outputLines := []string{}
	for i:= 0; foundItems < len(ts.workload); i++ {
		wl, exists := ts.workload[i]
		if (exists) {
			foundItems++
		}
		for resId, work := range wl {
			outputLines = append(outputLines, fmt.Sprintf("|%4d |%3d |%4.f |", i, resId, work))
		}
	}

	return fmt.Sprintf("%s\n%s", ts.task, strings.Join(outputLines, "\n"))
}

func (ts TaskSchedule) getStartInterval() int {
	intervals := []int{}
	for interval, _ := range ts.workload {
		intervals = append(intervals, interval)
	}
	sort.Ints(intervals)
	return intervals[0]
}

func (ts TaskSchedule) getEndInterval() int {
	intervals := []int{}
	for interval, _ := range ts.workload {
		intervals = append(intervals, interval)
	}
	sort.Ints(intervals)
	return intervals[len(intervals) - 1]
}
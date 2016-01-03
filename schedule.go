package main
import (
	"time"
	"fmt"
)

type Schedule struct {
	project   *Project
	tasks     map[int]*TaskSchedule
	resources map[int]*ResourceLoad
}



// Initializes schedule by adding tasks and resources to the struct
func (s *Schedule) init(p *Project) {
	s.project = p
	for i := 0; i < len(s.project.Tasks); i++ {
		s.tasks[s.project.Tasks[i].Id] = &TaskSchedule{
			&s.project.Tasks[i],
			false,
			s.project.Tasks[i].Duration,
			make(map[int]map[int]float64)}
	}
	for i := 0; i < len(s.project.Resources); i++ {
		s.resources[s.project.Resources[i].Id] = &ResourceLoad{
			&s.project.Resources[i],
			make(map[int]float64),
			make(map[int]map[int]float64)}
	}
}

func (s *Schedule) plan() {
	interval := 0
	planningLoop:
	for {
		//get tasks ready for execution
		tasks, hasTasks := s.getExecutableTasks()
		//  if there are no such tasks, we're done
		if !hasTasks {
			break
		}
		//  loop through the tasks
		for i := 0; i < len(tasks); i++ {
			//for each task get free resources for it
			resources, hasFree := s.getFreeResources(interval, (*tasks[i]).task.Resources)
			if !hasFree {
				continue
			}
			//loop through each resource and apply his work to the task
			for j := 0; j < len(resources); j++ {
				s.planTask(interval, tasks[i], resources[j])
				if tasks[i].finished {
					break
				}
			}
			//if resources are filled, increase interval counter and start over
			if _, hasFreeResources := s.getFreeResource(interval); !hasFreeResources {
				interval++
				continue planningLoop
			}
		}
		//need to check if tasks here are the same we have started from and if yes, means we have not assigned anything
		//thus need to increase interval counter
		updatedTasks, _ := s.getExecutableTasks()
		if sameTasks(tasks, updatedTasks) {
			interval++
		}
	}
}

func sameTasks(t1, t2 []*TaskSchedule) bool {
	if t2 == nil || t2 == nil {
		return false
	}
	if (len(t1) != len(t2)) {
		return false
	}
	for i := 0; i < len(t1); i++ {
		if t1[i].task.Id != t2[i].task.Id {
			return false
		}
	}
	return true
}

// Returns executable tasks. Executable are those which are not finished and have prerequisites finished
func (s *Schedule) getExecutableTasks() ([]*TaskSchedule, bool) {
	tasks := []*TaskSchedule{}
	for _, ts := range s.tasks {
		if !ts.finished && ts.prerequisitesFinished(s) {
			tasks = append(tasks, ts)
		}
	}
	//todo: sort by "flow' based on DueDate and Duration. Need to pass "interval" for that?
	return tasks, len(tasks) > 0
}

func (s *Schedule) getFreeResource(interval int) (*ResourceLoad, bool) {
	for _, rl := range s.resources {
		if rl.free(interval) > 0 {
			return rl, true
		}
	}
	return nil, false
}

// Returns slice of ResourceLoad pointers, a list of free resources for a given interval among the given resource list
func (s *Schedule) getFreeResources(interval int, requiredResources []int) ([]*ResourceLoad, bool) {
	resources := []*ResourceLoad{}
	for _, resourceId := range requiredResources {
		if s.resources[resourceId].free(interval) > 0 {
			resources = append(resources, s.resources[resourceId])
		}
	}
	return resources, len(resources) > 0
}

func (s *Schedule) planTask(interval int, ts *TaskSchedule, rl *ResourceLoad) {
	var plannedLoad float64 = 0
	if (ts.workRemained > rl.free(interval)) {
		plannedLoad = rl.free(interval)
	} else {
		plannedLoad = ts.workRemained
	}
	ts.work(interval, rl.resource.Id, plannedLoad)
	rl.work(interval, ts.task.Id, plannedLoad)
}

func (s *Schedule) export() ProjectExport {
	startTime := s.project.getStartTime()
	pe := ProjectExport{
		s.project.Name,
		fmt.Sprintf("/Date(%d)/", startTime.Unix() * 1000),
		make([]TaskExport, 0)}
	var te TaskExport
	for _, ts := range s.tasks {
		te = TaskExport{}
		te.Id = ts.task.Id
		te.Desc = ts.task.Name
		te.StartTime = startTime.Add(time.Duration(time.Hour * time.Duration(24 * ts.getStartInterval())))
		te.From = fmt.Sprintf("/Date(%d)/", te.StartTime.Unix() * 1000)
		te.EndTime = startTime.Add(time.Duration(time.Hour * time.Duration(24 * ts.getEndInterval()) + time.Hour * 8))
		te.To = fmt.Sprintf("/Date(%d)/", te.EndTime.Unix() * 1000)
		te.Deps = ts.task.BlockedBy
		pe.Tasks = append(pe.Tasks, te)
	}
	return pe
}

type ProjectExport struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	Tasks     []TaskExport `json:"tasks"`
}

type TaskExport struct {
	Id        int `json:"id"`
	Desc      string `json:"desc"`
	StartTime time.Time `json:"-"`
	From      string `json:"from"`
	EndTime   time.Time `json:"-"`
	To        string `json:"to"`
	Deps      []int  `json:"dep"`
}
package main
import (
	"sort"
	"time"
)

type Project struct {
	Name      string
	StartDate string
	Tasks     []Task
	Resources []Resource
}

type ByStartDate []Task

func (a ByStartDate) Len() int { return len(a) }
func (a ByStartDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByStartDate) Less(i, j int) bool { return a[i].startTime < a[j].startTime }

func (project *Project) sortTasks() {
	sort.Sort(ByStartDate(project.Tasks))
}

func (project *Project) resolveDeps() bool {
	resolved := true
	for i := 0; i < len(project.Tasks) && resolved; i++ {
		resolved = resolved && project.Tasks[i].buildPrerequisites(project)
	}
	return resolved
}

func (p Project) getStartTime() time.Time {
	tm, err := time.Parse(time.RFC3339, p.StartDate)
	if err != nil {
		tm = time.Now()
	}
	return tm
}
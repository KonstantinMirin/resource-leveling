package main
import "fmt"

type Task struct {
	Id                   int
	Name                 string
	Duration             float64
	Resources            []int
	BlockedBy            []int
//	DueDate              string

	startTime            float64
	endTime              float64

	prerequisites        []Task
	dependents           []Task
	resolvedDependencies bool
}

func (t *Task) adjustStart(startTime float64) {
	t.startTime = startTime;
	t.endTime = t.startTime + t.Duration
	if (len(t.dependents) > 0) {
		for i:= 0; i < len(t.dependents); i++ {
			t.dependents[i].adjustStart(t.endTime)
		}
	}
}

func (t *Task) buildPrerequisites(project *Project) bool {
	if len(t.BlockedBy) == 0 {
		t.resolvedDependencies = true
		return true
	}
	dependenciesFound := 0
	for i := 0; i < len(project.Tasks); i++ {
		for j := 0; j < len(t.BlockedBy); j++ {
			if (project.Tasks[i].Id == t.BlockedBy[j]) {
				t.prerequisites = append(t.prerequisites, project.Tasks[i])
				project.Tasks[i].dependents = append(project.Tasks[i].dependents, *t)
				dependenciesFound++
			}
		}
	}
	t.resolvedDependencies = dependenciesFound == len(t.BlockedBy)
	return t.resolvedDependencies
}

func (t Task) getPrerequisites(project *Project) *[]Task {
	if !t.resolvedDependencies {
		t.buildPrerequisites(project)
	}
	return &t.prerequisites
}

func (t Task) String() string {
	return fmt.Sprintf("%s: %.fh", t.Name, t.Duration)
}


func (t Task) assignableTo(resourceId int) bool {
	for i := 0; i < len(t.Resources); i++ {
		if resourceId == t.Resources[i] {
			return true
		}
	}
	return false
}

func (t Task) hasPrerequisites() bool {
	return len(t.prerequisites) > 0
}
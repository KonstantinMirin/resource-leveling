package main

import (
	"encoding/json"
	"flag"
//	"github.com/davecgh/go-spew/spew"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "project.json", "Path to project config file")
	flag.Parse()

	project := Project{}
	config, err := readConfig(configPath)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(config, &project); err != nil {
		panic(err)
	}
	project.resolveDeps()
	schedule := Schedule{&project, make(map[int]*TaskSchedule), make(map[int]*ResourceLoad)}
	schedule.init(&project)
	schedule.plan()

	printSchedule(schedule)
}

func printProject(p Project) {
	p.sortTasks()
	for i := 0; i < len(p.Tasks); i++ {
		fmt.Println(p.Tasks[i])
	}
}

func printSchedule(s Schedule) {
	exportObject := s.export()
	exportData, err := json.MarshalIndent(exportObject.Tasks, "", "    ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("public/out.json", exportData, os.FileMode(os.ModePerm))

	fmt.Println(string(exportData))
}
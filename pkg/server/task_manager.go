package server

import (
	"github.com/vmware/arachne/pkg/astrolabe"
	"sync"
	"time"
)

type TaskManager struct {
	tasks map[astrolabe.TaskID]astrolabe.Task
	mutex sync.RWMutex

	// For the clean up routine
	keepRunning bool
}

func NewTaskManager() TaskManager {
	newTM := TaskManager {
		keepRunning: true,
	}
	go newTM.cleanUpLoop()
	return newTM
}

func (this *TaskManager) ListTasks() []astrolabe.TaskID {
	this.mutex.RLock()
	defer this.mutex.Unlock()
	retTasks := make([]astrolabe.TaskID, len(this.tasks))
	curTaskNum := 0
	for curTask := range this.tasks {
		retTasks[curTaskNum] = curTask
		curTaskNum++
	}
	return retTasks
}

func (this *TaskManager) AddTask(addTask astrolabe.Task) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.tasks[addTask.GetID()] = addTask
}

func (this *TaskManager) RetrieveTask(taskID astrolabe.TaskID) (retTask astrolabe.Task, ok bool) {
	this.mutex.RLock()
	defer this.mutex.Unlock()
	retTask, ok =  this.tasks[taskID]
	return
}

func (this * TaskManager) cleanUpLoop() {
	for this.keepRunning {
		this.cleanUp()
		time.Sleep(time.Minute)
	}
}

func (this *TaskManager) cleanUp() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for id, task := range this.tasks {
		if time.Now().Sub(task.GetFinishedTime()) > 3600 * time.Second {
			delete(this.tasks, id)
		}
	}
}
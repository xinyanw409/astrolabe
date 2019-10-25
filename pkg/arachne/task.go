package arachne

import (
	"github.com/vmware/arachne/gen/models"
	"time"
)

type TaskID struct {
	id string
}

type TaskStatus int
const (
	Running TaskStatus = iota
	Success
	Failed
	Cancelled
)

type TaskResult interface {
	GetTaskStatus() TaskStatus
	GetResults() interface {}
	GetResultsString() string
}

func (t TaskStatus) String() string {
	return [...]string{"running", "success", "failed", "cancelled"}[t]
}

type Task interface {
	GetID() TaskID
	GetDetails() string
	GetFinishedTime() time.Time
	GetStartedTime() time.Time
	GetProgress() float64
	GetResult() TaskResult
	GetModelTaskInfo() (models.TaskInfo)
	Cancel() error
}


type GenericTask struct {
	// Fields are public so we don't need setters from users of this structure
	ID TaskID
	Completed bool
	TaskStatus TaskStatus
	Details string
	StartedTime, FinishedTime time.Time
	Progress float64
}

func (this GenericTask) GetTaskStatus() TaskStatus {
	return this.TaskStatus
}

func (this GenericTask) GetTaskID() TaskID {
	return this.ID
}

func (this GenericTask) GetDetails() string {
	return this.Details
}

func (this GenericTask)  GetFinishedTime() time.Time {
	return this.FinishedTime
}

func (this GenericTask)  GetStartedTime() time.Time {
	return this.StartedTime
}

func (this GenericTask)  GetProgress() float64 {
	return this.Progress
}

func (this GenericTask)  GetStatus() TaskStatus {
	return this.TaskStatus
}

func (this GenericTask)  GetModelTaskInfo() (models.TaskInfo) {
	startedTimeStr := this.StartedTime.Format(time.RFC3339)
	return models.TaskInfo{
		Completed:    &this.Completed,
		Details:      "",
		FinishedTime: this.FinishedTime.Format(time.RFC3339),
		ID:           "",
		Progress:     &this.Progress,
		StartedTime:  &startedTimeStr,
		Status:       nil,
	}
}

func (this GenericTask)  Cancel() error {
	return nil
}
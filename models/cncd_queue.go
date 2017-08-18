package models

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/cncd/queue"
)

// Task defines scheduled pipeline Task.
type Task struct {
	ID     string            ``
	Data   []byte            ``
	Labels map[string]string `xorm:"JSON"`
}

func (t Task) TableName() string {
	return "cncd_task"
}

func TaskInsert(task *Task) error {
	_, err := x.Insert(task)
	return err
}

func TaskDelete(id string) error {
	_, err := x.Delete(&Task{
		ID :id,
	})
	return err
}

func TaskList() ([]*Task, error) {
	tasks := make([]*Task, 0)

	if err := x.Find(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// WithTaskStore returns a queue that is backed by the TaskStore. This
// ensures the task Queue can be restored when the system starts.
func WithTaskStore(q queue.Queue) queue.Queue {
	tasks, _ := TaskList()
	for _, task := range tasks {
		q.Push(context.Background(), &queue.Task{
			ID:     task.ID,
			Data:   task.Data,
			Labels: task.Labels,
		})
	}
	return &persistentQueue{q}
}

type persistentQueue struct {
	queue.Queue
}

// Push pushes an task to the tail of this queue.
func (q *persistentQueue) Push(c context.Context, task *queue.Task) error {
	TaskInsert(&Task{
		ID:     task.ID,
		Data:   task.Data,
		Labels: task.Labels,
	})
	err := q.Queue.Push(c, task)
	if err != nil {
		TaskDelete(task.ID)
	}
	return err
}

// Poll retrieves and removes a task head of this queue.
func (q *persistentQueue) Poll(c context.Context, f queue.Filter) (*queue.Task, error) {
	task, err := q.Queue.Poll(c, f)
	if task != nil {
		logrus.Debugf("pull queue item: %s: remove from backup", task.ID)
		if derr := TaskDelete(task.ID); derr != nil {
			logrus.Errorf("pull queue item: %s: failed to remove from backup: %s", task.ID, derr)
		} else {
			logrus.Debugf("pull queue item: %s: successfully removed from backup", task.ID)
		}
	}
	return task, err
}

// Evict removes a pending task from the queue.
func (q *persistentQueue) Evict(c context.Context, id string) error {
	err := q.Queue.Evict(c, id)
	if err == nil {
		TaskDelete(id)
	}
	return err
}

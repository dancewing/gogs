package repo

import (
	"fmt"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"
)

func PipelineCallback(c *context.APIContext, callback form.PipelineCallback) {

	task, err := models.GetPipelineHookTaskByUUID(callback.DeliveryID)

	if err != nil {
		c.JSON(200, form.ResponseError{
			Message: fmt.Sprintf("Task with UUID : %s not exist", callback.DeliveryID),
		})
	}

	job, err := models.GetJobByTask(task.ID)

	if err != nil {
		c.JSON(200, form.ResponseError{
			Message: fmt.Sprintf("Job with TaskID : %s not exist", task.ID),
		})
	}

	nextJob, err := job.FindNextJob()
	if err != nil {

	} else {
		err = nextJob.LoadAttributes()

		if err != nil {
			c.JSON(200, form.ResponseError{
				Message: fmt.Sprintf("Error Load Job (%d) Attributes", nextJob.ID),
			})
		}

		err = models.PreparePipelineNextHook(nextJob.Repository, nextJob)

		if err != nil {
			c.JSON(200, form.ResponseError{
				Message: fmt.Sprintf("Error PreparePipelineNextHook (%d) ", nextJob.ID),
			})
		}

	}

	c.JSONSuccess(form.Response{})
}

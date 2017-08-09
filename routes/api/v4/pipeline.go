package v4

import (
	"github.com/gogits/gogs/pkg/context"
	"github.com/gogits/gogs/pkg/form"

	"github.com/gogits/gogs/models"
	log "gopkg.in/clog.v1"
)

func PipelineCallback(c *context.APIContext, callback form.PipelineCallback) {

	log.Info("call back %v", callback)

	if callback.DeliveryID != "" {

		pipeline, err := models.GetPipelineByDeliveryID(callback.DeliveryID)
		if err != nil {

			pipeline = &models.Pipeline{
				Status:       "running",
				RepositoryID: c.Repo.Repository.ID,
				DeliveryUUID: callback.DeliveryID,
			}

			_, err = models.CreatePipeline(pipeline)

			if err != nil {
				c.JSON(500, form.ResponseError{
					Success: false,
					Message: "error create pipeline",
				})
				return
			}
		}

		if callback.Step != "" {
			job, err := models.GetJobByStep(callback.Step, pipeline.ID)

			if err != nil {
				job = &models.Job{
					PipelineID: pipeline.ID,
					Stage:      callback.Step,
					Status:     callback.State,
				}

				if callback.State == "running" {

				}

				_, err = models.CreateJob(job)

				if err != nil {
					c.JSON(500, form.ResponseError{
						Success: false,
						Message: "error CreateJob",
					})
					return
				}
			}

			job.Status = callback.State

			err = models.UpdateJob(job)

			if err != nil {
				c.JSON(500, form.ResponseError{
					Success: false,
					Message: "error UpdateJob",
				})
				return
			}

		}

		models.UpdatePipelineStatus(pipeline)

	}
	//
	//task, err := models.GetPipelineHookTaskByUUID(callback.DeliveryID)
	//
	//if err != nil {
	//	c.JSON(200, form.ResponseError{
	//		Message: fmt.Sprintf("Task with UUID : %s not exist", callback.DeliveryID),
	//	})
	//}
	//
	//job, err := models.GetJobByTask(task.ID)
	//
	//if err != nil {
	//	c.JSON(200, form.ResponseError{
	//		Message: fmt.Sprintf("Job with TaskID : %s not exist", task.ID),
	//	})
	//}
	//
	//nextJob, err := job.FindNextJob()
	//if err != nil {
	//
	//} else {
	//	err = nextJob.LoadAttributes()
	//
	//	if err != nil {
	//		c.JSON(200, form.ResponseError{
	//			Message: fmt.Sprintf("Error Load Job (%d) Attributes", nextJob.ID),
	//		})
	//	}
	//
	//	err = models.PreparePipelineNextHook(nextJob.Repository, nextJob)
	//
	//	if err != nil {
	//		c.JSON(200, form.ResponseError{
	//			Message: fmt.Sprintf("Error PreparePipelineNextHook (%d) ", nextJob.ID),
	//		})
	//	}
	//
	//}

	c.JSONSuccess(form.Response{})
}

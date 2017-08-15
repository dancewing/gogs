package form

import (
	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type NewPipelineHook struct {
	ContentType int `binding:"Required"`
	Secret      string
	JenkinsUser string `binding:"Required"`
	JenkinsHost string `binding:"Required"`
	Webhook
}

func (f *NewPipelineHook) Validate(ctx *macaron.Context, errs binding.Errors) binding.Errors {
	return validate(errs, ctx.Data, f, ctx.Locale)
}

type PipelineCallback struct {
	Number            int64  `json:"number"`
	Result            string `json:"result"`
	TimeInMillis      int64  `json:"timeInMillis"`
	StartTimeInMillis int64  `json:"startTimeInMillis"`
	Description       string `json:"description"`
	DisplayName       string `json:"displayName"`
	DeliveryID        string `json:"deliveryID"`
	State             string `json:"state"`  ///pending, running, canceled, success, failed
	Step              string `json:"step"`
}

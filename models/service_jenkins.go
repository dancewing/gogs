package models

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gogits/gogs/pkg/httplib"
	"github.com/gogits/gogs/pkg/setting"

	log "gopkg.in/clog.v1"
)

type JenkinsServiceConfigLoad struct {
	JenkinsHost  string `json:"jenkins_host"`
	JenkinsUser  string `json:"jenkins_user"`
	JenkinsToken string `json:"jenkins_token"`
	*ServiceConfigLoad
}

func (load *JenkinsServiceConfigLoad) Deliver(t *ServiceTask) error {

	t.IsDelivered = true

	timeout := time.Duration(setting.Webhook.DeliverTimeout) * time.Second

	req := httplib.Post(t.URL).SetTimeout(timeout, timeout).
		Header("X-Gogs-Delivery", t.UUID).
		Header("X-Gogs-Signature", t.Signature).
		Header("X-Gogs-Event", string(t.EventType)).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: setting.Webhook.SkipTLSVerify})

	//Header("X-Gogs-Callback", string(t.CallbackURL)).


	req = req.Header("Content-Type", "application/json").Body(t.PayloadContent)

	//switch t.ContentType {
	//case JSON:
	//	req = req.Header("Content-Type", "application/json").Body(t.PayloadContent)
	//case FORM:
	//	req.Param("payload", t.PayloadContent)
	//}

	// Record delivery information.
	t.RequestInfo = &HookRequest{
		Headers: map[string]string{},
	}
	for k, vals := range req.Headers() {
		t.RequestInfo.Headers[k] = strings.Join(vals, ",")
	}

	t.ResponseInfo = &HookResponse{
		Headers: map[string]string{},
	}

	defer func() {
		t.Delivered = time.Now().UnixNano()
		if t.IsSucceed {
			log.Trace("Hook delivered: %s", t.UUID)
		} else {
			log.Trace("Hook delivery failed: %s", t.UUID)
		}

		// Update webhook last delivery status.
		w, err := GetServiceByID(t.ConfigID)
		if err != nil {
			log.Error(3, "GetWebhookByID: %v", err)
			return
		}
		if t.IsSucceed {
			w.LastStatus = SERVICE_STATUS_SUCCEED
		} else {
			w.LastStatus = SERVICE_STATUS_FAILED
		}
		if err = UpdateServiceConfig(w); err != nil {
			log.Error(3, "UpdateWebhook: %v", err)
			return
		}
	}()

	resp, err := req.Response()
	if err != nil {
		t.ResponseInfo.Body = fmt.Sprintf("Delivery: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Status code is 20x can be seen as succeed.
	t.IsSucceed = resp.StatusCode/100 == 2
	t.ResponseInfo.Status = resp.StatusCode
	for k, vals := range resp.Header {
		t.ResponseInfo.Headers[k] = strings.Join(vals, ",")
	}

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.ResponseInfo.Body = fmt.Sprintf("read body: %s", err)
		return err
	}
	t.ResponseInfo.Body = string(p)

	return nil
}

func (config *ServiceConfig) ToJenkinsServiceConfigEdit() *JenkinsServiceConfigLoad {
	return ToJenkinsServiceConfigLoad(config)
}

func ToJenkinsServiceConfigLoad(config *ServiceConfig) *JenkinsServiceConfigLoad {
	if config == nil {
		return &JenkinsServiceConfigLoad{ServiceConfigLoad: &ServiceConfigLoad{
			HookEvent: &HookEvent{},
		}}
	}

	jenkinsServiceConfig := &JenkinsServiceConfigLoad{}
	json.Unmarshal([]byte(config.ConfigContent), jenkinsServiceConfig)

	jenkinsServiceConfig.ReadEvent()

	jenkinsServiceConfig.ServiceConfigLoad.IsActive = config.IsActive
	jenkinsServiceConfig.ServiceConfigLoad.ConfigID = config.ID
	jenkinsServiceConfig.ServiceConfigLoad.RepoID = config.RepoID
	jenkinsServiceConfig.ServiceConfigLoad.OrgID = config.OrgID

	return jenkinsServiceConfig
}

var _ ServiceDelivery = new(JenkinsServiceConfigLoad)

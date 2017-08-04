// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"encoding/json"
	"time"

	"github.com/go-xorm/xorm"
	log "gopkg.in/clog.v1"

	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/httplib"
	"github.com/gogits/gogs/pkg/setting"
	"github.com/gogits/gogs/pkg/sync"
	gouuid "github.com/satori/go.uuid"
)

var PipelineHookQueue = sync.NewUniqueQueue(setting.Webhook.QueueLength)

// PipelineHook represents a hook object trigger by web hook.
type PipelineHook struct {
	ID          int64
	RepoID      int64
	OrgID       int64
	ContentType HookContentType
	Secret      string `xorm:"TEXT"`
	Events      string `xorm:"TEXT"`
	*HookEvent  `xorm:"-"`
	IsSSL       bool `xorm:"is_ssl"`
	IsActive    bool
	Meta        string     `xorm:"TEXT"` // store hook-specific attributes
	LastStatus  HookStatus // Last delivery status

	Created     time.Time `xorm:"-"`
	CreatedUnix int64
	Updated     time.Time `xorm:"-"`
	UpdatedUnix int64
}

func (w *PipelineHook) BeforeInsert() {
	w.CreatedUnix = time.Now().Unix()
	w.UpdatedUnix = w.CreatedUnix
}

func (w *PipelineHook) BeforeUpdate() {
	w.UpdatedUnix = time.Now().Unix()
}

func (w *PipelineHook) AfterSet(colName string, _ xorm.Cell) {
	var err error
	switch colName {
	case "events":
		w.HookEvent = &HookEvent{}
		if err = json.Unmarshal([]byte(w.Events), w.HookEvent); err != nil {
			log.Error(3, "Unmarshal [%d]: %v", w.ID, err)
		}
	case "created_unix":
		w.Created = time.Unix(w.CreatedUnix, 0).Local()
	case "updated_unix":
		w.Updated = time.Unix(w.UpdatedUnix, 0).Local()
	}
}

// CreateWebhook creates a new web hook.
func CreatePipelineHook(w *PipelineHook) error {
	_, err := x.Insert(w)
	return err
}

// UpdateWebhook updates information of webhook.
func UpdatePipelineHook(w *PipelineHook) error {
	_, err := x.Id(w.ID).AllCols().Update(w)
	return err
}

// History returns history of webhook by given conditions.
func (w *PipelineHook) History(page int) ([]*PipelineHookTask, error) {
	return PipelineHookTasks(w.ID, page)
}

// HookTasks returns a list of hook tasks by given conditions.
func PipelineHookTasks(hookID int64, page int) ([]*PipelineHookTask, error) {
	tasks := make([]*PipelineHookTask, 0, setting.Webhook.PagingNum)
	return tasks, x.Limit(setting.Webhook.PagingNum, (page-1)*setting.Webhook.PagingNum).Where("hook_id=?", hookID).Desc("id").Find(&tasks)
}

// UpdateEvent handles conversion from HookEvent to Events.
func (w *PipelineHook) UpdateEvent() error {
	data, err := json.Marshal(w.HookEvent)
	w.Events = string(data)
	return err
}

// HasCreateEvent returns true if hook enabled create event.
func (w *PipelineHook) HasCreateEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Create)
}

// HasDeleteEvent returns true if hook enabled delete event.
func (w *PipelineHook) HasDeleteEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Delete)
}

// HasForkEvent returns true if hook enabled fork event.
func (w *PipelineHook) HasForkEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Fork)
}

// HasPushEvent returns true if hook enabled push event.
func (w *PipelineHook) HasPushEvent() bool {
	return w.PushOnly || w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Push)
}

// HasIssuesEvent returns true if hook enabled issues event.
func (w *PipelineHook) HasIssuesEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Issues)
}

// HasIssueCommentEvent returns true if hook enabled issue comment event.
func (w *PipelineHook) HasIssueCommentEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.IssueComment)
}

// HasPullRequestEvent returns true if hook enabled pull request event.
func (w *PipelineHook) HasPullRequestEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.PullRequest)
}

// HasReleaseEvent returns true if hook enabled release event.
func (w *PipelineHook) HasReleaseEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Release)
}

// HookTask represents a hook task.
type PipelineHookTask struct {
	ID              int64
	RepoID          int64 `xorm:"INDEX"`
	HookID          int64
	UUID            string
	URL             string `xorm:"TEXT"`
	Signature       string `xorm:"TEXT"`
	api.Payloader   `xorm:"-"`
	PayloadContent  string `xorm:"TEXT"`
	ContentType     HookContentType
	EventType       HookEventType
	IsSSL           bool
	IsDelivered     bool
	Delivered       int64
	DeliveredString string `xorm:"-"`

	CallbackURL string `xorm:"TEXT"`

	// History info.
	IsSucceed       bool
	RequestContent  string        `xorm:"TEXT"`
	RequestInfo     *HookRequest  `xorm:"-"`
	ResponseContent string        `xorm:"TEXT"`
	ResponseInfo    *HookResponse `xorm:"-"`
}

func (t *PipelineHookTask) BeforeUpdate() {
	if t.RequestInfo != nil {
		t.RequestContent = t.MarshalJSON(t.RequestInfo)
	}
	if t.ResponseInfo != nil {
		t.ResponseContent = t.MarshalJSON(t.ResponseInfo)
	}
}

func (t *PipelineHookTask) AfterSet(colName string, _ xorm.Cell) {
	var err error
	switch colName {
	case "delivered":
		t.DeliveredString = time.Unix(0, t.Delivered).Format("2006-01-02 15:04:05 MST")

	case "request_content":
		if len(t.RequestContent) == 0 {
			return
		}

		t.RequestInfo = &HookRequest{}
		if err = json.Unmarshal([]byte(t.RequestContent), t.RequestInfo); err != nil {
			log.Error(3, "Unmarshal[%d]: %v", t.ID, err)
		}

	case "response_content":
		if len(t.ResponseContent) == 0 {
			return
		}

		t.ResponseInfo = &HookResponse{}
		if err = json.Unmarshal([]byte(t.ResponseContent), t.ResponseInfo); err != nil {
			log.Error(3, "Unmarshal [%d]: %v", t.ID, err)
		}
	}
}

func (t *PipelineHookTask) MarshalJSON(v interface{}) string {
	p, err := json.Marshal(v)
	if err != nil {
		log.Error(3, "Marshal [%d]: %v", t.ID, err)
	}
	return string(p)
}

// getActiveWebhooksByRepoID returns all active webhooks of repository.
func getActivePipelineHooksByRepoID(e Engine, repoID int64) ([]*PipelineHook, error) {
	webhooks := make([]*PipelineHook, 0, 5)
	return webhooks, e.Where("repo_id = ?", repoID).And("is_active = ?", true).Find(&webhooks)
}

// getActiveWebhooksByOrgID returns all active webhooks for an organization.
func getActivePipelineHooksByOrgID(e Engine, orgID int64) ([]*PipelineHook, error) {
	ws := make([]*PipelineHook, 0, 3)
	return ws, e.Where("org_id=?", orgID).And("is_active=?", true).Find(&ws)
}

// GetWebhookOfRepoByID returns webhook of repository by given ID.
func GetPipelineHookOfRepoByID(repoID, id int64) (*PipelineHook, error) {
	return getPipelineHook(&PipelineHook{
		//	ID:     id,
		RepoID: repoID,
	})
}

// GetWebhookByOrgID returns webhook of organization by given ID.
func GetPipelineHookByOrgID(orgID, id int64) (*PipelineHook, error) {
	return getPipelineHook(&PipelineHook{
		//	ID:    id,
		OrgID: orgID,
	})
}

// getWebhook uses argument bean as query condition,
// ID must be specified and do not assign unnecessary fields.
func getPipelineHook(bean *PipelineHook) (*PipelineHook, error) {
	has, err := x.Get(bean)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.WebhookNotExist{bean.ID}
	}
	return bean, nil
}

// GetWebhookByID returns webhook by given ID.
// Use this function with caution of accessing unauthorized webhook,
// which means should only be used in non-user interactive functions.
func GetPipelineHookByID(id int64) (*PipelineHook, error) {
	return getPipelineHook(&PipelineHook{
		ID: id,
	})
}

// prepareHookTasks adds list of webhooks to task queue.
func preparePipelineHookTasks(e Engine, repo *Repository, event HookEventType, p api.Payloader, webhooks []*PipelineHook, job *Job) (err error) {
	if len(webhooks) == 0 {
		return nil
	}

	var payloader api.Payloader
	for _, w := range webhooks {
		switch event {
		case HOOK_EVENT_CREATE:
			if !w.HasCreateEvent() {
				continue
			}
		case HOOK_EVENT_DELETE:
			if !w.HasDeleteEvent() {
				continue
			}
		case HOOK_EVENT_FORK:
			if !w.HasForkEvent() {
				continue
			}
		case HOOK_EVENT_PUSH:
			if !w.HasPushEvent() {
				continue
			}
		case HOOK_EVENT_ISSUES:
			if !w.HasIssuesEvent() {
				continue
			}
		case HOOK_EVENT_ISSUE_COMMENT:
			if !w.HasIssueCommentEvent() {
				continue
			}
		case HOOK_EVENT_PULL_REQUEST:
			if !w.HasPullRequestEvent() {
				continue
			}
		case HOOK_EVENT_RELEASE:
			if !w.HasReleaseEvent() {
				continue
			}
		}

		payloader = p

		var signature string
		if len(w.Secret) > 0 {
			data, err := payloader.JSONPayload()
			if err != nil {
				log.Error(2, "prepareWebhooks.JSONPayload: %v", err)
			}
			sig := hmac.New(sha256.New, []byte(w.Secret))
			sig.Write(data)
			signature = hex.EncodeToString(sig.Sum(nil))
		}

		callbackURL := setting.AppURL + "api/v1/pipeline/callback"

		hookTask := &PipelineHookTask{
			RepoID:      repo.ID,
			HookID:      w.ID,
			URL:         job.JobURL,
			Signature:   signature,
			Payloader:   payloader,
			ContentType: w.ContentType,
			EventType:   event,
			IsSSL:       w.IsSSL,
			CallbackURL: callbackURL,
		}

		if err = createPipelineHookTask(e, hookTask); err != nil {
			return fmt.Errorf("createHookTask: %v", err)
		}

		job.HookTaskID = hookTask.ID
		job.DeliveryUUID = hookTask.UUID
		err = UpdateJob(job)
		if err != nil {
			return fmt.Errorf("Error UpdateJob: %v", err)
		}

	}

	// It's safe to fail when the whole function is called during hook execution
	// because resource released after exit. Also, there is no process started to
	// consume this input during hook execution.
	go PipelineHookQueue.Add(repo.ID)
	return nil
}

// createHookTask creates a new hook task,
// it handles conversion from Payload to PayloadContent.
func createPipelineHookTask(e Engine, t *PipelineHookTask) error {
	data, err := t.Payloader.JSONPayload()
	if err != nil {
		return err
	}
	t.UUID = gouuid.NewV4().String()
	t.PayloadContent = string(data)
	_, err = e.Insert(t)
	return err
}

// UpdateHookTask updates information of hook task.
func UpdatePipelineHookTask(t *PipelineHookTask) error {
	_, err := x.Id(t.ID).AllCols().Update(t)
	return err
}

func getPipelineHookTaskByID(e Engine, id int64) (*PipelineHookTask, error) {
	m := &PipelineHookTask{
		ID: id,
	}
	has, err := e.Get(m)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrPipelineNotExist{id}
	}
	return m, nil
}

func GetPipelineHookTask(taskID int64) (*PipelineHookTask, error) {
	return getPipelineHookTaskByID(x, taskID)
}

func GetPipelineHookTaskByUUID(uuid string) (*PipelineHookTask, error) {
	task := &PipelineHookTask{
		UUID: uuid,
	}
	has, err := x.Get(task)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrPipelineHookTaskNotExist{0, uuid}
	}
	return task, nil
}

func preparePipelineHooks(e Engine, repo *Repository, event HookEventType, p api.Payloader) error {
	webhooks, err := getActivePipelineHooksByRepoID(e, repo.ID)
	if err != nil {
		return fmt.Errorf("getActiveWebhooksByRepoID [%d]: %v", repo.ID, err)
	}

	// check if repo belongs to org and append additional webhooks
	if repo.mustOwner(e).IsOrganization() {
		// get hooks for org
		orgws, err := getActivePipelineHooksByOrgID(e, repo.OwnerID)
		if err != nil {
			return fmt.Errorf("getActiveWebhooksByOrgID [%d]: %v", repo.OwnerID, err)
		}
		webhooks = append(webhooks, orgws...)
	}

	def, err := GetCIFileFromGit(repo.Owner, repo)
	if err != nil {
		return fmt.Errorf("Error GetCIFileFromGit: %v", err)
	}

	job, err := CreatePipeline(e, def, repo)
	if err != nil {
		return fmt.Errorf("Error CreatePipeline: %v", err)
	}

	return preparePipelineHookTasks(e, repo, event, p, webhooks, job)
}

func PreparePipelineNextHook(repo *Repository, job *Job) error {

	webhooks, err := getActivePipelineHooksByRepoID(x, repo.ID)
	if err != nil {
		return fmt.Errorf("getActiveWebhooksByRepoID [%d]: %v", repo.ID, err)
	}

	// check if repo belongs to org and append additional webhooks
	if repo.mustOwner(x).IsOrganization() {
		// get hooks for org
		orgws, err := getActivePipelineHooksByOrgID(x, repo.OwnerID)
		if err != nil {
			return fmt.Errorf("getActiveWebhooksByOrgID [%d]: %v", repo.OwnerID, err)
		}
		webhooks = append(webhooks, orgws...)
	}

	p := &HookPayload{
		Repo: repo.APIFormat(nil),
	}

	return preparePipelineHookTasks(x, repo, "hook", p, webhooks, job)
}

func (t *PipelineHookTask) deliver() {
	t.IsDelivered = true

	timeout := time.Duration(setting.Webhook.DeliverTimeout) * time.Second
	req := httplib.Post(t.URL).SetTimeout(timeout, timeout).
		Header("X-Gogs-Delivery", t.UUID).
		Header("X-Gogs-Signature", t.Signature).
		Header("X-Gogs-Event", string(t.EventType)).
		Header("X-Gogs-Callback", string(t.CallbackURL)).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: setting.Webhook.SkipTLSVerify})

	switch t.ContentType {
	case JSON:
		req = req.Header("Content-Type", "application/json").Body(t.PayloadContent)
	case FORM:
		req.Param("payload", t.PayloadContent)
	}

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
		w, err := GetPipelineHookByID(t.HookID)
		if err != nil {
			log.Error(3, "GetWebhookByID: %v", err)
			return
		}
		if t.IsSucceed {
			w.LastStatus = HOOK_STATUS_SUCCEED
		} else {
			w.LastStatus = HOOK_STATUS_FAILED
		}
		if err = UpdatePipelineHook(w); err != nil {
			log.Error(3, "UpdateWebhook: %v", err)
			return
		}
	}()

	resp, err := req.Response()
	if err != nil {
		t.ResponseInfo.Body = fmt.Sprintf("Delivery: %v", err)
		return
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
		return
	}
	t.ResponseInfo.Body = string(p)
}

func DeliverPipelineHooks() {
	tasks := make([]*PipelineHookTask, 0, 10)
	x.Where("is_delivered = ?", false).Iterate(new(PipelineHookTask),
		func(idx int, bean interface{}) error {
			t := bean.(*PipelineHookTask)
			t.deliver()
			tasks = append(tasks, t)
			return nil
		})

	// Update hook task status.
	for _, t := range tasks {
		if err := UpdatePipelineHookTask(t); err != nil {
			log.Error(4, "UpdateHookTask [%d]: %v", t.ID, err)
		}
	}

	// Start listening on new hook requests.
	for repoID := range PipelineHookQueue.Queue() {
		log.Trace("DeliverHooks [repo_id: %v]", repoID)
		PipelineHookQueue.Remove(repoID)

		tasks = make([]*PipelineHookTask, 0, 5)
		if err := x.Where("repo_id = ?", repoID).And("is_delivered = ?", false).Find(&tasks); err != nil {
			log.Error(4, "Get repository [%s] hook tasks: %v", repoID, err)
			continue
		}
		for _, t := range tasks {
			t.deliver()
			if err := UpdatePipelineHookTask(t); err != nil {
				log.Error(4, "UpdateHookTask [%d]: %v", t.ID, err)
				continue
			}
		}
	}
}

func InitPipelineDeliverHooks() {
	go DeliverPipelineHooks()
}

type ErrPipelineHookNotExist struct {
	ID int64
}

func IsErrPipelineHookNotExist(err error) bool {
	_, ok := err.(ErrPipelineHookNotExist)
	return ok
}

func (err ErrPipelineHookNotExist) Error() string {
	return fmt.Sprintf("PipelineHook does not exist [id: %d ]", err.ID)
}

type ErrPipelineHookTaskNotExist struct {
	ID   int64
	UUID string
}

func IsErrPipelineHookTaskNotExist(err error) bool {
	_, ok := err.(ErrPipelineHookNotExist)
	return ok
}

func (err ErrPipelineHookTaskNotExist) Error() string {
	return fmt.Sprintf("PipelineHookTask does not exist [id: %d, UUID : %s]", err.ID, err.UUID)
}

type HookPayload struct {
	Repo *api.Repository `json:"repository"`
	api.Payloader
}

func (p *HookPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

var _ api.Payloader = &HookPayload{}

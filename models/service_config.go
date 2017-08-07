package models

import (
	"time"

	"encoding/json"

	"fmt"

	"crypto/tls"
	"io/ioutil"
	"strings"

	"github.com/go-xorm/xorm"
	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/models/errors"
	"github.com/gogits/gogs/pkg/httplib"
	"github.com/gogits/gogs/pkg/setting"
	"github.com/gogits/gogs/pkg/sync"
	gouuid "github.com/satori/go.uuid"
	log "gopkg.in/clog.v1"
)

var ServiceQueue = sync.NewUniqueQueue(setting.Webhook.QueueLength)

type ServiceType int

const (
	JENKINS ServiceType = iota + 1
)

type ServiceStatus int

const (
	SERVICE_STATUS_NONE = iota
	SERVICE_STATUS_SUCCEED
	SERVICE_STATUS_FAILED
)

var serviceTypes = map[string]ServiceType{
	"jenkins": JENKINS,
}

// ToHookTaskType returns HookTaskType by given name.
func ToServiceType(name string) ServiceType {
	return serviceTypes[name]
}

func (t ServiceType) Name() string {
	switch t {
	case JENKINS:
		return "jenkins"
	}
	return ""
}

type ServiceConfig struct {
	ID     int64
	Type   ServiceType
	RepoID int64
	OrgID  int64

	Name string `xorm:"-"`

	ConfigContent string `xorm:"TEXT"`
	IsActive      bool
	Created       time.Time `xorm:"-"`
	CreatedUnix   int64
	Updated       time.Time `xorm:"-"`
	UpdatedUnix   int64

	LastStatus ServiceStatus

	*HookEvent `xorm:"-"`
}

func (sc *ServiceConfig) BeforeInsert() {
	sc.CreatedUnix = time.Now().Unix()
	sc.UpdatedUnix = sc.CreatedUnix
}

func (sc *ServiceConfig) BeforeUpdate() {
	sc.UpdatedUnix = time.Now().Unix()
}

func (sc *ServiceConfig) AfterSet(colName string, _ xorm.Cell) {

	var err error
	switch colName {
	case "created_unix":
		sc.Created = time.Unix(sc.CreatedUnix, 0).Local()
	case "updated_unix":
		sc.Updated = time.Unix(sc.UpdatedUnix, 0)
	case "type":
		sc.Name = sc.Type.Name()

	case "config_content":
		sc.HookEvent = &HookEvent{}
		if err = json.Unmarshal([]byte(sc.ConfigContent), sc.HookEvent); err != nil {
			log.Error(3, "Unmarshal [%d]: %v", sc.ID, err)
		}
	}
}

// HasCreateEvent returns true if hook enabled create event.
func (w *ServiceConfig) HasCreateEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Create)
}

// HasDeleteEvent returns true if hook enabled delete event.
func (w *ServiceConfig) HasDeleteEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Delete)
}

// HasForkEvent returns true if hook enabled fork event.
func (w *ServiceConfig) HasForkEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Fork)
}

// HasPushEvent returns true if hook enabled push event.
func (w *ServiceConfig) HasPushEvent() bool {
	return w.PushOnly || w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Push)
}

// HasIssuesEvent returns true if hook enabled issues event.
func (w *ServiceConfig) HasIssuesEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Issues)
}

// HasIssueCommentEvent returns true if hook enabled issue comment event.
func (w *ServiceConfig) HasIssueCommentEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.IssueComment)
}

// HasPullRequestEvent returns true if hook enabled pull request event.
func (w *ServiceConfig) HasPullRequestEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.PullRequest)
}

// HasReleaseEvent returns true if hook enabled release event.
func (w *ServiceConfig) HasReleaseEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Release)
}

func GetConfiguration(serviceType ServiceType, repositoryID int64) (*ServiceConfig, error) {
	return getConfiguration(&ServiceConfig{
		Type:   serviceType,
		RepoID: repositoryID,
	})
}

func CreateConfiguration(config *ServiceConfig, content interface{}) (*ServiceConfig, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return config, err
	}
	config.ConfigContent = string(data)
	_, err = x.Insert(config)
	return config, err
}

func UpdateConfiguration(config *ServiceConfig, content interface{}) (*ServiceConfig, error) {
	data, err := json.Marshal(content)
	if err != nil {
		return config, err
	}
	config.ConfigContent = string(data)
	_, err = x.AllCols().Update(config)
	return config, err
}

func getConfiguration(bean *ServiceConfig) (*ServiceConfig, error) {
	has, err := x.Get(bean)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.ConfigurationNotExist{Type: bean.Type.Name(), RepositoryID: bean.RepoID}
	}
	return bean, nil
}

func GetServicesByRepoID(repoID int64) ([]*ServiceConfig, error) {
	services := make([]*ServiceConfig, 0, 5)
	return services, x.Find(&services, &ServiceConfig{RepoID: repoID})
}

func GetAllServicesByRepoID(repoID int64) ([]*ServiceConfig, error) {
	services := make([]*ServiceConfig, 0, 5)
	err := x.Find(&services, &ServiceConfig{RepoID: repoID})
	if err != nil {
		return nil, err
	}

	allServices := make([]*ServiceConfig, len(serviceTypes))

	index := 0
	for _, i := range serviceTypes {
		found := false

		for _, s := range services {
			if i == s.Type {
				allServices[index] = s
				found = true
			}
		}
		if !found {
			srvCfg := ServiceConfig{Name: i.Name()}
			allServices[index] = &srvCfg
		}
		index++
	}
	return allServices, nil
}

type ServiceConfigLoad struct {
	Events       string `json:"-"`
	Create       bool   `json:"-"`
	Delete       bool   `json:"-"`
	Fork         bool   `json:"-"`
	Push         bool   `json:"-"`
	Issues       bool   `json:"-"`
	IssueComment bool   `json:"-"`
	PullRequest  bool   `json:"-"`
	Release      bool   `json:"-"`
	IsActive     bool   `json:"-"`

	*HookEvent
}

func (f ServiceConfigLoad) PushOnly() bool {
	return f.Events == "push_only"
}

func (f ServiceConfigLoad) SendEverything() bool {
	return f.Events == "send_everything"
}

func (f ServiceConfigLoad) ChooseEvents() bool {
	return f.Events == "choose_events"
}

func (config *ServiceConfigLoad) ReadEvent() {

	if config.HookEvent == nil {
		return
	}

	if config.HookEvent.PushOnly {
		config.Events = "push_only"
	}
	if config.HookEvent.SendEverything {
		config.Events = "send_everything"
	}
	if config.HookEvent.ChooseEvents {
		config.Events = "choose_events"
	}

	config.Create = config.HookEvent.HookEvents.Create
	config.Delete = config.HookEvent.HookEvents.Delete
	config.Fork = config.HookEvent.HookEvents.Fork
	config.Push = config.HookEvent.HookEvents.Push
	config.Issues = config.HookEvent.HookEvents.Issues
	config.IssueComment = config.HookEvent.HookEvents.IssueComment
	config.PullRequest = config.HookEvent.HookEvents.PullRequest
	config.Release = config.HookEvent.HookEvents.Release

}

func (config *ServiceConfigLoad) UpdateEvent() {

	if config.HookEvent == nil {
		config.HookEvent = &HookEvent{}
	}

	config.HookEvent.PushOnly = config.PushOnly()
	config.HookEvent.SendEverything = config.SendEverything()
	config.HookEvent.ChooseEvents = config.ChooseEvents()

	if config.ChooseEvents() {
		config.HookEvent.HookEvents.Create = config.Create
		config.HookEvent.HookEvents.Delete = config.Delete
		config.HookEvent.HookEvents.Fork = config.Fork
		config.HookEvent.HookEvents.Push = config.Push
		config.HookEvent.HookEvents.Issues = config.Issues
		config.HookEvent.HookEvents.IssueComment = config.IssueComment
		config.HookEvent.HookEvents.PullRequest = config.PullRequest
		config.HookEvent.HookEvents.Release = config.Release
	} else {
		config.HookEvent.HookEvents.Create = false
		config.HookEvent.HookEvents.Delete = false
		config.HookEvent.HookEvents.Fork = false
		config.HookEvent.HookEvents.Push = false
		config.HookEvent.HookEvents.Issues = false
		config.HookEvent.HookEvents.IssueComment = false
		config.HookEvent.HookEvents.PullRequest = false
		config.HookEvent.HookEvents.Release = false
	}
}

// HookTask represents a hook task.
type ServiceTask struct {
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

func (t *ServiceTask) BeforeUpdate() {
	if t.RequestInfo != nil {
		t.RequestContent = t.MarshalJSON(t.RequestInfo)
	}
	if t.ResponseInfo != nil {
		t.ResponseContent = t.MarshalJSON(t.ResponseInfo)
	}
}

func (t *ServiceTask) AfterSet(colName string, _ xorm.Cell) {
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

func (t *ServiceTask) MarshalJSON(v interface{}) string {
	p, err := json.Marshal(v)
	if err != nil {
		log.Error(3, "Marshal [%d]: %v", t.ID, err)
	}
	return string(p)
}

// getActiveWebhooksByRepoID returns all active webhooks of repository.
func getActiveServicesByRepoID(e Engine, repoID int64) ([]*ServiceConfig, error) {
	webhooks := make([]*ServiceConfig, 0, 5)
	return webhooks, e.Where("repo_id = ?", repoID).And("is_active = ?", true).Find(&webhooks)
}

// getActiveWebhooksByOrgID returns all active webhooks for an organization.
func getActiveServicesByOrgID(e Engine, orgID int64) ([]*ServiceConfig, error) {
	ws := make([]*ServiceConfig, 0, 3)
	return ws, e.Where("org_id=?", orgID).And("is_active=?", true).Find(&ws)
}

func prepareServices(e Engine, repo *Repository, event HookEventType, p api.Payloader) error {

	webhooks, err := getActiveServicesByRepoID(e, repo.ID)
	if err != nil {
		return fmt.Errorf("getActiveWebhooksByRepoID [%d]: %v", repo.ID, err)
	}

	// check if repo belongs to org and append additional webhooks
	if repo.mustOwner(e).IsOrganization() {
		// get hooks for org
		orgws, err := getActiveServicesByOrgID(e, repo.OwnerID)
		if err != nil {
			return fmt.Errorf("getActiveWebhooksByOrgID [%d]: %v", repo.OwnerID, err)
		}
		webhooks = append(webhooks, orgws...)
	}

	return prepareServiceTasks(e, repo, event, p, webhooks)
}

// prepareHookTasks adds list of webhooks to task queue.
func prepareServiceTasks(e Engine, repo *Repository, event HookEventType, p api.Payloader, webhooks []*ServiceConfig) (err error) {
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

		callbackURL := setting.AppURL + "api/v1/pipeline/callback"

		hookTask := &ServiceTask{
			RepoID:      repo.ID,
			HookID:      w.ID,
			Payloader:   payloader,
			EventType:   event,
			CallbackURL: callbackURL,
		}

		if err = createServiceTask(e, hookTask); err != nil {
			return fmt.Errorf("createHookTask: %v", err)
		}

	}

	// It's safe to fail when the whole function is called during hook execution
	// because resource released after exit. Also, there is no process started to
	// consume this input during hook execution.
	go ServiceQueue.Add(repo.ID)
	return nil
}

func createServiceTask(e Engine, t *ServiceTask) error {
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
func UpdateServiceTask(t *ServiceTask) error {
	_, err := x.Id(t.ID).AllCols().Update(t)
	return err
}

// UpdateHookTask updates information of hook task.
func UpdateServiceConfig(t *ServiceConfig) error {
	_, err := x.Id(t.ID).AllCols().Update(t)
	return err
}

// getWebhook uses argument bean as query condition,
// ID must be specified and do not assign unnecessary fields.
func getService(bean *ServiceConfig) (*ServiceConfig, error) {
	has, err := x.Get(bean)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("Service don't exist ")
	}
	return bean, nil
}

// GetWebhookByID returns webhook by given ID.
// Use this function with caution of accessing unauthorized webhook,
// which means should only be used in non-user interactive functions.
func GetServiceByID(id int64) (*ServiceConfig, error) {
	return getService(&ServiceConfig{
		ID: id,
	})
}

func (t *ServiceTask) deliver() {
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
		w, err := GetServiceByID(t.HookID)
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

func DeliverServices() {
	tasks := make([]*ServiceTask, 0, 10)
	x.Where("is_delivered = ?", false).Iterate(new(ServiceTask),
		func(idx int, bean interface{}) error {
			t := bean.(*ServiceTask)
			t.deliver()
			tasks = append(tasks, t)
			return nil
		})

	// Update hook task status.
	for _, t := range tasks {
		if err := UpdateServiceTask(t); err != nil {
			log.Error(4, "UpdateServiceTask [%d]: %v", t.ID, err)
		}
	}

	// Start listening on new hook requests.
	for repoID := range ServiceQueue.Queue() {
		log.Trace("ServiceTask [repo_id: %v]", repoID)
		ServiceQueue.Remove(repoID)

		tasks = make([]*ServiceTask, 0, 5)
		if err := x.Where("repo_id = ?", repoID).And("is_delivered = ?", false).Find(&tasks); err != nil {
			log.Error(4, "Get repository [%s] service tasks: %v", repoID, err)
			continue
		}
		for _, t := range tasks {
			t.deliver()
			if err := UpdateServiceTask(t); err != nil {
				log.Error(4, "UpdateServiceTask [%d]: %v", t.ID, err)
				continue
			}
		}
	}
}

func InitServices() {
	go DeliverServices()
}

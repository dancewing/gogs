package models

import (
	"time"

	"encoding/json"

	"fmt"

	"github.com/go-xorm/xorm"
	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/models/errors"
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

// HookTask represents a hook task.
type ServiceTask struct {
	ID     int64
	RepoID int64 `xorm:"INDEX"`

	ConfigID       int64
	Config         *ServiceConfig `xorm:"-"`
	UUID           string
	URL            string `xorm:"TEXT"`
	Signature      string `xorm:"TEXT"`
	api.Payloader  `xorm:"-"`
	PayloadContent string `xorm:"TEXT"`

	EventType      HookEventType
	ServiceType     ServiceType

	IsSSL           bool
	IsDelivered     bool
	Delivered       int64
	DeliveredString string `xorm:"-"`

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
		return fmt.Errorf("getActiveServicesByRepoID [%d]: %v", repo.ID, err)
	}

	// check if repo belongs to org and append additional webhooks
	if repo.mustOwner(e).IsOrganization() {
		// get hooks for org
		orgws, err := getActiveServicesByOrgID(e, repo.OwnerID)
		if err != nil {
			return fmt.Errorf("getActiveServicesByOrgID [%d]: %v", repo.OwnerID, err)
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

		hookTask := &ServiceTask{
			RepoID:      repo.ID,
			ConfigID:    w.ID,
			Payloader:   payloader,
			EventType:   event,
			ServiceType: w.Type,
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
	var delivery ServiceDelivery
	switch t.ServiceType {
	case JENKINS:
		delivery = ToJenkinsServiceConfigLoad(t.Config)
	}

	if delivery != nil {
		delivery.Deliver(t)
	}
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

type ServiceDelivery interface {
	Deliver(task *ServiceTask) error
}

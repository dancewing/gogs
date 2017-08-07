package models

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
	ConfigID     int64  `json:"-"`
	RepoID       int64  `json:"-"`
	OrgID        int64  `json:"-"`
	*HookEvent
}

//
//func (f ServiceConfigLoad) PushOnly() bool {
//	return f.Events == "push_only"
//}
//
//func (f ServiceConfigLoad) SendEverything() bool {
//	return f.Events == "send_everything"
//}
//
//func (f ServiceConfigLoad) ChooseEvents() bool {
//	return f.Events == "choose_events"
//}

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

	config.HookEvent.PushOnly = config.Events == "push_only"
	config.HookEvent.SendEverything = config.Events == "send_everything"
	config.HookEvent.ChooseEvents = config.Events == "choose_events"

	if config.ChooseEvents {
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

// HasCreateEvent returns true if hook enabled create event.
func (w *ServiceConfigLoad) HasCreateEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Create)
}

// HasDeleteEvent returns true if hook enabled delete event.
func (w *ServiceConfigLoad) HasDeleteEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Delete)
}

// HasForkEvent returns true if hook enabled fork event.
func (w *ServiceConfigLoad) HasForkEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Fork)
}

// HasPushEvent returns true if hook enabled push event.
func (w *ServiceConfigLoad) HasPushEvent() bool {
	return w.PushOnly || w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Push)
}

// HasIssuesEvent returns true if hook enabled issues event.
func (w *ServiceConfigLoad) HasIssuesEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Issues)
}

// HasIssueCommentEvent returns true if hook enabled issue comment event.
func (w *ServiceConfigLoad) HasIssueCommentEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.IssueComment)
}

// HasPullRequestEvent returns true if hook enabled pull request event.
func (w *ServiceConfigLoad) HasPullRequestEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.PullRequest)
}

// HasReleaseEvent returns true if hook enabled release event.
func (w *ServiceConfigLoad) HasReleaseEvent() bool {
	return w.SendEverything ||
		(w.ChooseEvents && w.HookEvents.Release)
}

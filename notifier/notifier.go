package notifier

type Notifier interface {
	Notify(notification *Notification) error
}

type Interview struct {
	Name        string
	Start       string
	Interviewer string
}

// Notification is passed to notifiers.
type Notification struct {
	Candidate  string
	Role       string
	Interviews []*Interview
}

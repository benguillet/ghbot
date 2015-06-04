package notifier

type Notifier interface {
	Notify(notification *Notification) error
}

// Notification is passed to notifiers.
type Notification struct {
	Candidate string
	// Interviewers []Interviewers
}

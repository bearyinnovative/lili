package notifier

type NotifierType interface {
	Notify(id, text string, media []string) error
}

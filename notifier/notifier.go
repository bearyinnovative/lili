package notifier

type NotifierType interface {
	Notify(text string, media []string) error
}

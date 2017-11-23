package notifier

type NotifierType interface {
	Notify(text string, images []string) error
}

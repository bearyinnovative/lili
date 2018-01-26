package telegram

type Notifier struct {
	Token string `yaml:"token"`
	// `@channel_name` or integer id as string: `-123456`
	ChatID                string `yaml:"chat_id"`
	ParseMode             string `yaml:"parse_mode,omitempty"`
	DisableNotification   bool   `yaml:"disable_notification"`
	DisableWebPagePreview bool   `yaml:"disable_web_page_preview"`
}

func (n *Notifier) Notify(id, text string, media []string) error {
	return n.notify(toPairs(text, media))
}

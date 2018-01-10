package telegram

type Notifier struct {
	Token string `yaml:"token"`
	// `@channel_name` or integer id as string: `-123456`
	ChatID    string `yaml:"chat_id"`
	ParseMode string `yaml:"parse_mode,omitempty"`
}

func (n *Notifier) Notify(id, text string, media []string) error {
	return n.notify(toPairs(text, media))
}

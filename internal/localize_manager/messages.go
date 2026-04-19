package localize_manager

func (l *Localizator) StartMessage() string {
	msg, _ := l.mapper["/start"]

	return msg
}

func (l *Localizator) SaveUserMessage() string {
	msg, _ := l.mapper["/save_user"]

	return msg
}

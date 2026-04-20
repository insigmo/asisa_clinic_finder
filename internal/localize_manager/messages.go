package localize_manager

func (l *Localizator) StartMessage() string {
	msg, _ := l.mapper["/start"]

	return msg
}

func (l *Localizator) SaveUserMessage() string {
	msg, _ := l.mapper["/save_user"]

	return msg
}

func (l *Localizator) ChangeCityMessage() string {
	msg, _ := l.mapper["/change_city"]

	return msg
}

func (l *Localizator) WrongCityMessage() string {
	msg, _ := l.mapper["/wrong_city"]

	return msg
}

func (l *Localizator) WrongDirection() string {
	msg, _ := l.mapper["/wrong_direction"]

	return msg
}

func (l *Localizator) Perhaps() string {
	msg, _ := l.mapper["/perhaps"]

	return msg
}

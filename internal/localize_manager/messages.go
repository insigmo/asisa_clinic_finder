package localize_manager

func (l *Localizator) StartMessage() string {
	return l.mapper["/start"]
}

func (l *Localizator) SaveUserMessage() string {
	return l.mapper["/save_user"]
}

func (l *Localizator) ChangeCityMessage() string {
	return l.mapper["/change_city"]
}

func (l *Localizator) WrongCityMessage() string {
	return l.mapper["/wrong_city"]
}

func (l *Localizator) WrongDirection() string {
	return l.mapper["/wrong_direction"]
}

func (l *Localizator) Perhaps() string {
	return l.mapper["/perhaps"]
}

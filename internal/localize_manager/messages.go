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

func (l *Localizator) AskDirectionMessage() string {
	return l.mapper["/ask_direction"]
}

func (l *Localizator) UnknownMessage() string {
	return l.mapper["/unknown_call"]
}
func (l *Localizator) ClinicsNotFound() string {
	return l.mapper["/clinics_not_found"]
}
func (l *Localizator) Address() string {
	return l.mapper["/address"]
}
func (l *Localizator) Phone() string {
	return l.mapper["/phone"]
}

func (l *Localizator) Clinic() string {
	return l.mapper["/clinic"]
}

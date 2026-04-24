package i18n

func (l *Localizer) StartMessage() string        { return l.mapper["/start"] }
func (l *Localizer) SaveUserMessage() string     { return l.mapper["/save_user"] }
func (l *Localizer) ChangeCityMessage() string   { return l.mapper["/change_city"] }
func (l *Localizer) WrongCityMessage() string    { return l.mapper["/wrong_city"] }
func (l *Localizer) WrongDirection() string      { return l.mapper["/wrong_direction"] }
func (l *Localizer) Perhaps() string             { return l.mapper["/perhaps"] }
func (l *Localizer) AskDirectionMessage() string { return l.mapper["/ask_direction"] }
func (l *Localizer) UnknownMessage() string      { return l.mapper["/unknown_call"] }
func (l *Localizer) ClinicsNotFound() string     { return l.mapper["/clinics_not_found"] }
func (l *Localizer) Clinic() string              { return l.mapper["/clinic"] }
func (l *Localizer) Address() string             { return l.mapper["/address"] }
func (l *Localizer) Phone() string               { return l.mapper["/phone"] }

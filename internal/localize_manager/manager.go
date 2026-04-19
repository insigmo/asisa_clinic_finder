package localize_manager

import "fmt"

type messageMapper map[string]string

var (
	rusMessages = messageMapper{
		"/start":     "Добро пожаловать. Пожалуйста введите город в котором вы хотите найти поликлинику",
		"/save_user": "Спасибо! Ваш город сохранен, если вы хотите изменить свой город, вы можете это сделать через меню!",
	}
	esMessages = messageMapper{
		"/start":     "Bienvenido/a. Por favor, introducer la ciudad en la que desea encontrar una clínica.",
		"/save_user": "¡Gracias! Tu ciudad ya está guardada. Si quieres cambiarla, puedes hacerlo desde el menú.",
	}
	enMessages = messageMapper{
		"/start":     "Welcome. Please enter the city in which you would like to find a clinic.",
		"/save_user": "Thank you! Your city has been saved. If you want to change your city, you can do so via the menu!",
	}
	allMessages = []messageMapper{rusMessages, esMessages, enMessages}
)

type Localizator struct {
	mapper messageMapper
}

func New(languageCode string) *Localizator {
	validateMessages()

	var mapper messageMapper

	switch languageCode {
	case "ru":
		mapper = rusMessages
	case "es":
		mapper = esMessages
	case "en":
		mapper = enMessages
	}

	return &Localizator{
		mapper: mapper,
	}
}

func validateMessages() {
	expectLength := len(allMessages[0])

	for _, msgs := range allMessages[1:] {
		currentLength := len(msgs)
		if expectLength != currentLength {
			panic(fmt.Sprintf("expect length %d, got %d. %v", expectLength, currentLength, msgs))
		}
	}
}

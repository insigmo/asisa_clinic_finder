package localize_manager

import "fmt"

type messageMapper map[string]string

// TODO Попросить поправить текста через ии с мед смайликами

var (
	rusMessages = messageMapper{
		"/start":           `Добро пожаловать\. Пожалуйста введите город в котором вы хотите найти поликлинику`,
		"/change_city":     `Пожалуйста, укажите город, который необходимо сохранить`,
		"/wrong_city":      `Извините, данный город не найден в Испании\. Повторите попытку с другим городом`,
		"/save_user":       `Спасибо\! Ваш город сохранен, если вы хотите изменить свой город, вы можете это сделать через меню\!`,
		"/wrong_direction": `Извините, но мы не нашли медицинское направление\.`,
		"/perhaps":         ` Возможно вы имели ввиду: `,
		"/ask_direction":   `Пожалуйста, укажите медицинское направление, по которому необходимо подобрать поликлинику`,
		"/unknown_call":    `Извините, неизвестная команда`,
	}
	esMessages = messageMapper{
		"/start":           `Bienvenido/a\. Por favor, introducer la ciudad en la que desea encontrar una clínica\.`,
		"/change_city":     `Por favor, indique la ciudad que desea guardar`,
		"/wrong_city":      `Lo sentimos, no se ha encontrado esa ciudad en España\. Inténtalo de nuevo con otra ciudad\.`,
		"/save_user":       `¡Gracias\! Tu ciudad ya está guardada\. Si quieres cambiarla, puedes hacerlo desde el menú\.`,
		"/wrong_direction": `Lo sentimos, pero no hemos encontrado ninguna derivación médica\.`,
		"/perhaps":         ` Quizás se refería a: `,
		"/ask_direction":   `Por favor, indique la especialidad médica para la que desea encontrar una clínica`,
		"/unknown_call":    `Perdón, comando desconocido`,
	}
	enMessages = messageMapper{
		"/start":           `Welcome\. Please enter the city in which you would like to find a clinic\.`,
		"/change_city":     `Please specify the city you would like to save`,
		"/wrong_city":      `Sorry, we couldn't find that city in Spain\. Please try again with a different city\.`,
		"/save_user":       `Thank you\! Your city has been saved\. If you want to change your city, you can do so via the menu\!`,
		"/wrong_direction": `Sorry, but we couldn't find a medical referral\.`,
		"/perhaps":         ` Perhaps you meant: `,
		"/ask_direction":   `Please specify the medical specialty for which you would like to find a clinic`,
		"/unknown_call":    `Sorry, unknown command`,
	}
	allMessages = []messageMapper{rusMessages, esMessages, enMessages}
)

type Localizator struct {
	mapper messageMapper
}

func New(languageCode string) *Localizator {
	validateMessages()

	mapper := enMessages

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

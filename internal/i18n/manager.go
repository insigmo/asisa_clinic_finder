package i18n

import (
	"fmt"
	"strings"
	"sync"
)

type messageMapper map[string]string

var (
	rusMessages = messageMapper{
		"/start":             `Добро пожаловать\. Пожалуйста, укажите город, в котором вы хотите найти медицинское учреждение\.`,
		"/change_city":       `Пожалуйста, укажите город, который необходимо сохранить\.`,
		"/wrong_city":        `К сожалению, указанный город не найден в Испании\. Пожалуйста, попробуйте указать другой город\.`,
		"/save_user":         `Город успешно сохранён\. Если вы хотите изменить его, воспользуйтесь соответствующим пунктом меню\.`,
		"/wrong_direction":   `К сожалению, указанное медицинское направление не найдено\.`,
		"/perhaps":           ` Возможно, вы имели в виду: `,
		"/ask_direction":     `Пожалуйста, укажите медицинское направление, по которому необходимо подобрать учреждение\.`,
		"/unknown_call":      `Неизвестная команда\. Пожалуйста, воспользуйтесь кнопками меню\.`,
		"/clinics_not_found": `К сожалению, в базе данных не найдено учреждений с указанной специальностью\. Попробуйте указать смежное направление — например, не андролог, а уролог\.`,
		"/clinic":            `Учреждение: `,
		"/address":           `Адрес: `,
		"/phone":             `Телефон: `,
	}
	esMessages = messageMapper{
		"/start":             `Bienvenido/a\. Por favor, indique la ciudad en la que desea encontrar un centro médico\.`,
		"/change_city":       `Por favor, indique la ciudad que desea guardar\.`,
		"/wrong_city":        `Lo sentimos, la ciudad indicada no ha sido encontrada en España\. Por favor, inténtelo de nuevo con otra ciudad\.`,
		"/save_user":         `La ciudad ha sido guardada correctamente\. Si desea modificarla, puede hacerlo desde el menú correspondiente\.`,
		"/wrong_direction":   `Lo sentimos, la especialidad médica indicada no ha sido encontrada\.`,
		"/perhaps":           ` Quizás se refería a: `,
		"/ask_direction":     `Por favor, indique la especialidad médica para la que desea encontrar un centro\.`,
		"/unknown_call":      `Comando desconocido\. Por favor, utilice los botones del menú\.`,
		"/clinics_not_found": `Lo sentimos, no se han encontrado centros con la especialidad indicada en nuestra base de datos\. Pruebe con una especialidad relacionada — por ejemplo, en lugar de andrólogo, urólogo\.`,
		"/clinic":            `Centro médico: `,
		"/address":           `Dirección: `,
		"/phone":             `Teléfono: `,
	}
	enMessages = messageMapper{
		"/start":             `Welcome\. Please enter the city in which you would like to find a medical centre\.`,
		"/change_city":       `Please specify the city you would like to save\.`,
		"/wrong_city":        `The city you entered was not found in Spain\. Please try again with a different city\.`,
		"/save_user":         `Your city has been saved successfully\. If you wish to change it, you may do so via the menu\.`,
		"/wrong_direction":   `The medical specialty you entered was not found\.`,
		"/perhaps":           ` Perhaps you meant: `,
		"/ask_direction":     `Please specify the medical specialty for which you would like to find a centre\.`,
		"/unknown_call":      `Unknown command\. Please use the buttons below\.`,
		"/clinics_not_found": `No centres with the specified specialty were found in our database\. Please try a related specialty — for example, instead of an andrologist, try a urologist\.`,
		"/clinic":            `Medical centre: `,
		"/address":           `Address: `,
		"/phone":             `Phone: `,
	}

	once sync.Once
)

// Localizer предоставляет локализованные сообщения.
type Localizer struct{ mapper messageMapper }

// New возвращает Localizer для заданного языкового кода (ru/es/en, по умолчанию en).
func New(languageCode string) *Localizer {
	once.Do(func() {
		want := len(enMessages)
		for _, m := range []messageMapper{rusMessages, esMessages} {
			if len(m) != want {
				panic(fmt.Sprintf("i18n: map length mismatch: want %d, got %d", want, len(m)))
			}
		}
	})
	switch normalize(languageCode) {
	case "ru":
		return &Localizer{mapper: rusMessages}
	case "es":
		return &Localizer{mapper: esMessages}
	default:
		return &Localizer{mapper: enMessages}
	}
}

func normalize(code string) string {
	code = strings.TrimSpace(strings.ToLower(code))
	switch {
	case strings.HasPrefix(code, "ru"):
		return "ru"
	case strings.HasPrefix(code, "es"):
		return "es"
	default:
		return "en"
	}
}

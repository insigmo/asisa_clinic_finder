package clinic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/insigmo/asisa_clinic_finder/internal/constants"
	"github.com/insigmo/asisa_clinic_finder/internal/db"
	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
	"github.com/insigmo/asisa_clinic_finder/internal/localize_manager"
)

const (
	maxPages           = 5
	countClinicsOnPage = 6
	maxClinics         = maxPages * countClinicsOnPage
)

func Search(params *local_models.BaseParams, dbManager *db.Manager, user *db.User, direction string) (string, error) {
	client := helpers.NewHTTPManager()
	places, err := client.FetchPlaces(params.Ctx, user.City)
	if err != nil {
		return "", fmt.Errorf("fetch places: %w", err)
	}
	if len(places) == 0 {
		return "", fmt.Errorf("town %q not found", user.City)
	}

	provinceID, err := client.FetchProvinceID(params.Ctx, places[0].PlaceID)
	if err != nil {
		return "", fmt.Errorf("fetch province id: %w", err)
	}
	allClinics := make([]local_models.Clinic, 0, maxClinics)

	for i := 1; i <= maxPages; i++ {
		data, err := client.FetchClinics(params.Ctx, provinceID, direction, i)
		if err != nil {
			return "", fmt.Errorf("fetch clinics: %w", err)
		}
		clinics, err := client.ParseHTML(data, direction)
		if err != nil {
			return "", fmt.Errorf("parse html: %w", err)
		}
		if len(clinics) == 0 {
			break
		}
		allClinics = append(allClinics, clinics...)
	}

	if len(allClinics) == 0 {
		return "", nil
	}
	allClinics = sortClinicsByPostalCode(allClinics, params, dbManager, user)
	return prepareResult(allClinics, user.LanguageCode), nil
}

func escapeMarkdownV2(text string) string {
	var builder strings.Builder
	builder.Grow(len(text) + len(text)/10) // небольшой запас

	for _, r := range text {
		switch r {
		case '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!':
			builder.WriteRune('\\')
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func sortClinicsByPostalCode(clinics []local_models.Clinic, params *local_models.BaseParams, dbManager *db.Manager, user *db.User) []local_models.Clinic {
	localizator := localize_manager.New(user.LanguageCode)
	postalCodes, err := dbManager.FindCity(params.Ctx, user.City)
	if err != nil {
		params.Log.Error(err.Error())
		return clinics
	}

	if len(postalCodes) == 0 {
		if err = helpers.SendMessage(params, localizator.WrongCityMessage()); err != nil {
			params.Log.Error(err.Error())
		}

		helpers.SetUserState(params, constants.StateIdle)
		return clinics
	}

	postalCode := postalCodes[len(postalCodes)/2]
	return sortClinics(clinics, postalCode)
}

func sortClinics(clinics []local_models.Clinic, postalCode int) []local_models.Clinic {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	sort.Slice(clinics, func(i, j int) bool {
		return abs(clinics[i].PostalCode-postalCode) < abs(clinics[j].PostalCode-postalCode)
	})
	return clinics
}

func prepareResult(clinics []local_models.Clinic, languageCode string) string {
	localizator := localize_manager.New(languageCode)
	if len(clinics) == 0 {
		return localizator.ClinicsNotFound()
	}

	var builder strings.Builder

	// "Заголовок" — жирный текст + эмодзи (настоящих заголовков в Telegram нет)
	direction := escapeMarkdownV2(strings.ToTitle(strings.ToLower(clinics[0].Direction)))
	builder.WriteString("🏥 *")
	builder.WriteString(direction)
	builder.WriteString("*\n")

	for i, c := range clinics {
		// Номер + название клиники жирным
		name := escapeMarkdownV2(strings.ToTitle(strings.ToLower(c.Name)))
		builder.WriteString(fmt.Sprintf("*%d\\. %s`%s`*\n", i+1, localizator.Clinic(), name))

		address := escapeMarkdownV2(strings.ToTitle(strings.ToLower(c.Address)))
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%s`%s`\n", localizator.Address(), address))
		builder.WriteString(fmt.Sprintf("%s`+34%s`\n", localizator.Phone(), c.PhoneNumber))
		builder.WriteString("\n\n")
	}

	return builder.String()
}

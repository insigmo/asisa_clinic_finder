package clinic

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func Search(ctx context.Context, city, direction string, postalCode int) (string, error) {
	client := helpers.NewHTTPManager()
	// TODO добавить пагинацию до 5 страниц
	places, err := client.FetchPlaces(ctx, city)
	if err != nil {
		return "", fmt.Errorf("fetch places: %w", err)
	}
	if len(places) == 0 {
		return "", fmt.Errorf("town %q not found", city)
	}

	provinceID, err := client.FetchProvinceID(ctx, places[0].PlaceID)
	if err != nil {
		return "", fmt.Errorf("fetch province id: %w", err)
	}

	data, err := client.FetchClinics(ctx, provinceID, direction)
	if err != nil {
		return "", fmt.Errorf("fetch clinics: %w", err)
	}

	clinics, err := client.ParseHTML(data, direction)
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}
	sortClinics(clinics, postalCode)
	return prepareResult(clinics), nil
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

func prepareResult(clinics []local_models.Clinic) string {
	if len(clinics) == 0 {
		// TODO поправить на мультиязычный вариант
		return "_No clinics found_"
	}

	var builder strings.Builder

	// "Заголовок" — жирный текст + эмодзи (настоящих заголовков в Telegram нет)
	direction := escapeMarkdownV2(strings.ToTitle(strings.ToLower(clinics[0].Direction)))
	builder.WriteString("🏥 *")
	builder.WriteString(direction)
	builder.WriteString("*\n")
	builder.WriteString(fmt.Sprintf("_Found %d clinic\\(s\\)_\n\n", len(clinics)))

	// TODO Добавить на 1 страницу 1 сообщение
	// TODO Добавить возможность копировать адрес
	for i, c := range clinics {
		// Номер + название клиники жирным
		name := escapeMarkdownV2(strings.ToTitle(strings.ToLower(c.Name)))
		builder.WriteString(fmt.Sprintf("*%d\\. %s*\n", i+1, name))

		address := escapeMarkdownV2(strings.ToTitle(strings.ToLower(c.Address)))
		builder.WriteString(fmt.Sprintf("📍 %s\n", address))
		builder.WriteString(fmt.Sprintf("📞 `+34%s`\n", c.PhoneNumber))
		builder.WriteString(fmt.Sprintf("🏷️ `%d`\n", c.PostalCode))

		if c.Distance > 0 {
			builder.WriteString(fmt.Sprintf("📏 %d m\n", c.Distance))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

func sortClinics(clinics []local_models.Clinic, postalCode int) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	sort.Slice(clinics, func(i, j int) bool {
		return abs(clinics[i].PostalCode-postalCode) < abs(clinics[j].PostalCode-postalCode)
	})
}

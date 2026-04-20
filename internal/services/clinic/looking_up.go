package clinic

import (
	"fmt"
	"sort"
	"strings"

	"github.com/insigmo/asisa_clinic_finder/internal/helpers"
	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

func Search(city, direction string, postalCode int) string {
	client := helpers.NewHttpManager()

	places, err := client.FetchPlaces(city)
	if err != nil {
		panic(err)
	}
	if len(places) == 0 {
		panic(fmt.Errorf("town '%s' not found", city))
	}

	provinceID, err := client.FetchProvinceID(places[0].PlaceId)
	if err != nil {
		panic(err)
	}

	data, err := client.FetchClinics(provinceID, direction)
	if err != nil {
		panic(err)
	}

	clinics, err := client.ParseHtml(data, direction)
	if err != nil {
		panic(err)
	}
	sortClinics(clinics, postalCode)
	return prepareResult(clinics)
}

func prepareResult(clinics []local_models.Clinic) string {
	var builder strings.Builder

	for _, clinic := range clinics {
		builder.WriteString(clinic.Direction + "\n")
		builder.WriteString(clinic.Name + "\n")
		builder.WriteString(clinic.Address + "\n")
		builder.WriteString(clinic.PhoneNumber + "\n")
		builder.WriteString(fmt.Sprintf("%d\n", clinic.PostalCode))
		builder.WriteString("\n")
	}

	return builder.String()
}

func sortClinics(clinics []local_models.Clinic, postalCode int) {
	sort.Slice(clinics, func(i, j int) bool {
		first := clinics[i].PostalCode - postalCode
		if first < 0 {
			first = -first
		}
		second := clinics[j].PostalCode - postalCode
		if second < 0 {
			second = -second
		}

		return first < second
	})
}

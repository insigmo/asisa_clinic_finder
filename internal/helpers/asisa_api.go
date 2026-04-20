package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

const (
	baseUrl = "https://www.asisa.es"
)

var postalCodePattern, _ = regexp.Compile(`\d{5}`)

func (h *HttpManager) FetchClinics(provinceID int, direction string) (string, error) {
	fetchUrl := fmt.Sprintf("%s/cuadro-medico/resultados-cuadro-medico", baseUrl)
	query := url.Values{}
	query.Set("networkId", "1")
	query.Set("specialityType", "1")
	query.Set("specialityName", strings.ToUpper(direction))
	query.Set("provinceId", strconv.Itoa(provinceID))

	body, err := h.SendGetRequest(fetchUrl, query.Encode())
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (h *HttpManager) FetchProvinceID(placeID string) (int, error) {
	fetchUrl := fmt.Sprintf("%s/bin/wasisa/gmaps-service", baseUrl)
	query := url.Values{}
	query.Set("id", placeID)
	body, err := h.SendGetRequest(fetchUrl, query.Encode())
	if err != nil {
		return 0, err
	}
	var province local_models.Province
	err = json.Unmarshal(body, &province)
	if err != nil {
		return 0, err
	}
	return province.ProvinceID, nil
}

func (h *HttpManager) FetchPlaces(town string) ([]local_models.Place, error) {
	fetchUrl := fmt.Sprintf("%s/bin/wasisa/autocomplete-addresses", baseUrl)

	query := url.Values{}
	query.Set("q", strings.ToLower(town))

	body, err := h.SendGetRequest(fetchUrl, query.Encode())
	if err != nil {
		return nil, err
	}
	places := make([]local_models.Place, 0, 3)

	err = json.Unmarshal(body, &places)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return places, nil
}

func (h *HttpManager) SendGetRequest(url string, values string) ([]byte, error) {
	fullUrl := fmt.Sprintf("%s?%s", url, values)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}
	return body, nil
}

func (h *HttpManager) ParseHtml(data, direction string) ([]local_models.Clinic, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	clinics := make([]local_models.Clinic, 0, 10)
	doc.Find("div.cmp-medical-picture-result__info-container").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__contact-data--name").Text())
		address := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container--address").Text())
		number := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__location--phone").Text())
		postal, _ := strconv.Atoi(postalCodePattern.FindString(address))
		clinics = append(clinics, local_models.Clinic{
			Name:        strings.ToLower(name),
			Direction:   direction,
			Address:     strings.ToLower(address),
			PhoneNumber: number,
			PostalCode:  postal,
		})
	})

	return clinics, nil
}

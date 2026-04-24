package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/insigmo/asisa_clinic_finder/internal/model"
)

const (
	baseURL   = "https://www.asisa.es"
	maxPlaces = 10
)

// postalCodePattern — компилируем один раз, переиспользуем при каждом парсинге.
var postalCodePattern = regexp.MustCompile(`\d{5}`)

func (h *HTTPManager) FetchClinics(ctx context.Context, provinceID int, direction string, page int) (string, error) {
	q := urlpkg.Values{}
	q.Set("networkId", "1")
	q.Set("specialityType", "1")
	q.Set("specialityName", strings.ToUpper(direction))
	q.Set("provinceId", strconv.Itoa(provinceID))
	q.Set("page", strconv.Itoa(page))
	body, err := h.sendGet(ctx, baseURL+"/cuadro-medico/resultados-cuadro-medico", q.Encode())
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (h *HTTPManager) FetchProvinceID(ctx context.Context, placeID string) (int, error) {
	q := urlpkg.Values{}
	q.Set("id", placeID)
	body, err := h.sendGet(ctx, baseURL+"/bin/wasisa/gmaps-service", q.Encode())
	if err != nil {
		return 0, err
	}
	var province model.Province
	if err := json.Unmarshal(body, &province); err != nil {
		return 0, fmt.Errorf("unmarshal province: %w", err)
	}
	return province.ProvinceID, nil
}

func (h *HTTPManager) FetchPlaces(ctx context.Context, town string) ([]model.Place, error) {
	q := urlpkg.Values{}
	q.Set("q", strings.ToLower(town))
	body, err := h.sendGet(ctx, baseURL+"/bin/wasisa/autocomplete-addresses", q.Encode())
	if err != nil {
		return nil, err
	}
	places := make([]model.Place, 0, maxPlaces)
	if err := json.Unmarshal(body, &places); err != nil {
		return nil, fmt.Errorf("unmarshal places: %w", err)
	}
	return places, nil
}

func (h *HTTPManager) sendGet(ctx context.Context, rawURL, values string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL+"?"+values, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (h *HTTPManager) ParseHTML(data, direction string) ([]model.Clinic, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}
	clinics := make([]model.Clinic, 0, maxPlaces)
	doc.Find("div.cmp-medical-picture-result__info-container").Each(func(_ int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__contact-data--name").Text())
		address := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container--address").Text())
		phone := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__location--phone").Text())
		postal, _ := strconv.Atoi(postalCodePattern.FindString(address))
		clinics = append(clinics, model.Clinic{
			Name:        strings.ToLower(name),
			Direction:   direction,
			Address:     strings.ToLower(address),
			PhoneNumber: phone,
			PostalCode:  postal,
		})
	})
	return clinics, nil
}

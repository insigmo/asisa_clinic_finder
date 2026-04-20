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

	"github.com/insigmo/asisa_clinic_finder/internal/local_models"
)

const baseURL = "https://www.asisa.es"

var postalCodePattern = regexp.MustCompile(`\d{5}`)

func (h *HTTPManager) FetchClinics(ctx context.Context, provinceID int, direction string) (string, error) {
	fetchURL := baseURL + "/cuadro-medico/resultados-cuadro-medico"
	q := urlpkg.Values{}
	q.Set("networkId", "1")
	q.Set("specialityType", "1")
	q.Set("specialityName", strings.ToUpper(direction))
	q.Set("provinceId", strconv.Itoa(provinceID))

	body, err := h.sendGet(ctx, fetchURL, q.Encode())
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (h *HTTPManager) FetchProvinceID(ctx context.Context, placeID string) (int, error) {
	fetchURL := baseURL + "/bin/wasisa/gmaps-service"
	q := urlpkg.Values{}
	q.Set("id", placeID)

	body, err := h.sendGet(ctx, fetchURL, q.Encode())
	if err != nil {
		return 0, err
	}
	var province local_models.Province
	if err := json.Unmarshal(body, &province); err != nil {
		return 0, fmt.Errorf("unmarshal province: %w", err)
	}
	return province.ProvinceID, nil
}

func (h *HTTPManager) FetchPlaces(ctx context.Context, town string) ([]local_models.Place, error) {
	fetchURL := baseURL + "/bin/wasisa/autocomplete-addresses"
	q := urlpkg.Values{}
	q.Set("q", strings.ToLower(town))

	body, err := h.sendGet(ctx, fetchURL, q.Encode())
	if err != nil {
		return nil, err
	}
	places := make([]local_models.Place, 0, 3)
	if err := json.Unmarshal(body, &places); err != nil {
		return nil, fmt.Errorf("unmarshal places: %w", err)
	}
	return places, nil
}

func (h *HTTPManager) sendGet(ctx context.Context, rawURL, values string) ([]byte, error) {
	fullURL := rawURL + "?" + values
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (h *HTTPManager) ParseHTML(data, direction string) ([]local_models.Clinic, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}
	clinics := make([]local_models.Clinic, 0, 10)
	doc.Find("div.cmp-medical-picture-result__info-container").Each(func(_ int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__contact-data--name").Text())
		address := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container--address").Text())
		phone := strings.TrimSpace(s.Find(".cmp-medical-picture-result__info-container__location--phone").Text())
		postal, _ := strconv.Atoi(postalCodePattern.FindString(address))
		clinics = append(clinics, local_models.Clinic{
			Name:        strings.ToLower(name),
			Direction:   direction,
			Address:     strings.ToLower(address),
			PhoneNumber: phone,
			PostalCode:  postal,
		})
	})
	return clinics, nil
}

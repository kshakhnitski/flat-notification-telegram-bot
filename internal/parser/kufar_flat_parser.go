package parser

import (
	"flat_bot/internal/model"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type KufarFlatParser struct {
	url string
}

func NewKufarFlatParser(url string) KufarFlatParser {
	return KufarFlatParser{url: url}
}

func (p KufarFlatParser) Parse() ([]model.Flat, error) {
	resp, err := http.Get(p.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var flats []model.Flat
	doc.Find("section").Each(func(index int, item *goquery.Selection) {
		link := extractLink(item)
		id, err := extractIDFromLink(link)
		if err != nil {
			return
		}

		params := extractParameters(item)
		flat := model.Flat{
			ID:          id,
			Parameters:  params.Summary,
			Rooms:       params.Rooms,
			AreaInSqM:   params.Area,
			Floor:       params.Floor,
			TotalFloors: params.TotalFloors,
			Address:     extractAddress(item),
			Description: extractDescription(item),
			Metro:       extractMetro(item),
			Link:        link,
			PriceInUsd:  extractPriceUSD(item),
			PriceInByn:  extractPriceBYN(item),
			Source:      "Kufar",
		}

		flats = append(flats, flat)
	})
	return flats, nil
}

func extractLink(item *goquery.Selection) string {
	link, _ := item.Find(".styles_wrapper__Q06m9").Attr("href")
	return link

}

type ApartmentParameters struct {
	Summary     string
	Rooms       int
	Area        *float64
	Floor       *int
	TotalFloors *int
}

func extractMetro(item *goquery.Selection) string {
	return item.Find(".styles_wrapper__HKXX4 span").Text()
}

func extractParameters(item *goquery.Selection) ApartmentParameters {
	parametersStr := item.Find(".styles_parameters__7zKlL").Text()

	parsedParams := ApartmentParameters{
		Summary:     parametersStr,
		Rooms:       0,
		Area:        nil,
		Floor:       nil,
		TotalFloors: nil,
	}

	roomsRegex := regexp.MustCompile(`(\d+)\s*комн\.`)
	if roomsMatch := roomsRegex.FindStringSubmatch(parametersStr); roomsMatch != nil {
		parsedParams.Rooms, _ = strconv.Atoi(roomsMatch[1])
	}

	areaRegex := regexp.MustCompile(`(\d+(\.\d+)?)\s*м²`)
	if areaMatch := areaRegex.FindStringSubmatch(parametersStr); areaMatch != nil {
		area, _ := strconv.ParseFloat(areaMatch[1], 64)
		parsedParams.Area = &area
	}

	floorRegex := regexp.MustCompile(`этаж\s*(\d+)(?:\s*из\s*(\d+))?`)
	if floorMatch := floorRegex.FindStringSubmatch(parametersStr); floorMatch != nil {
		floor, _ := strconv.Atoi(floorMatch[1])
		parsedParams.Floor = &floor
		if len(floorMatch) > 2 && floorMatch[2] != "" {
			totalFloors, _ := strconv.Atoi(floorMatch[2])
			parsedParams.TotalFloors = &totalFloors
		}
	}

	return parsedParams
}

func extractAddress(item *goquery.Selection) string {
	return item.Find(".styles_address__l6Qe_").Text()
}

func extractPriceBYN(item *goquery.Selection) float64 {
	text := item.Find(".styles_price__byr__lLSfd").Text()
	re := regexp.MustCompile(`\d{1,3}(?:\s?\d{3})*(?:\.\d+)?`)
	numberStr := re.FindString(text)
	numberStr = strings.Replace(numberStr, " ", "", -1)
	result, _ := strconv.ParseFloat(numberStr, 64)
	return result
}

func extractPriceUSD(item *goquery.Selection) float64 {
	text := item.Find(".styles_price__usd__HpXMa").Text()
	re := regexp.MustCompile(`\d{1,3}(?:\s?\d{3})*(?:\.\d+)?`)
	numberStr := re.FindString(text)
	numberStr = strings.Replace(numberStr, " ", "", -1)
	result, _ := strconv.ParseFloat(numberStr, 64)
	return result
}

func extractDescription(item *goquery.Selection) string {
	return item.Find(".styles_body__5BrnC, .styles_body__r33c8").Text()
}

func extractIDFromLink(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	pathSegments := strings.Split(parsedURL.Path, "/")
	id := pathSegments[len(pathSegments)-1]
	return id, nil
}

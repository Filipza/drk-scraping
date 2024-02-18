package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/playwright-community/playwright-go"
)

func main() {
	// Initialize Page
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	url := buildUrl(22041)

	if _, err = page.Goto(url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	getEntries(page)

	// Click next page button, wait for domcontentloaded and get entries again until no "next" btn found
	if ok, err := page.Locator(".next").IsVisible(); err == nil && ok {
		err = page.Locator(".next").Click()
		if err != nil {
			log.Fatalf("could not click next: %v", err)
		}

		page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateDomcontentloaded})

		getEntries(page)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func buildUrl(plz int) string {
	baseURL := "https://www.brk.de/angebote/kleidercontainer-suchergebnis.html?tx_drkclothescontainersearch_clothescontainersearch%5B__referrer%5D%5B%40extension%5D=DrkClothescontainersearch&tx_drkclothescontainersearch_clothescontainersearch%5B__referrer%5D%5B%40controller%5D=Clothescontainer&tx_drkclothescontainersearch_clothescontainersearch%5B__referrer%5D%5B%40action%5D=form&tx_drkclothescontainersearch_clothescontainersearch%5B__referrer%5D%5Barguments%5D=YTowOnt950e60808c4c0be500ae7f67c799028d239d52c53&tx_drkclothescontainersearch_clothescontainersearch%5B__referrer%5D%5B%40request%5D=%7B%22%40extension%22%3A%22DrkClothescontainersearch%22%2C%22%40controller%22%3A%22Clothescontainer%22%2C%22%40action%22%3A%22form%22%7D39a7d1855af4e35ea33af2efb91730f8b2579827&tx_drkclothescontainersearch_clothescontainersearch%5B__trustedProperties%5D=%7B%22sSearchstring%22%3A1%2C%22submit%22%3A1%7D55a7b1dc7819f186138d2feeec028f046320677c&tx_drkclothescontainersearch_clothescontainersearch%5BsSearchstring%5D="
	suffix := "&tx_drkclothescontainersearch_clothescontainersearch%5Bsubmit%5D=suchen"
	url := baseURL + strconv.Itoa(plz) + suffix

	return url
}

func getEntries(page playwright.Page) {
	entries, err := page.Locator(".t-medium-20.columns").All()
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	for _, v := range entries {
		mapUrl, err := v.Locator("a").GetAttribute("href")
		if err != nil {
			log.Fatalf("could not get href: %v", err)
		}
		getUrlCoordinates(mapUrl)

		textContent, err := v.InnerText()
		if err != nil {
			log.Fatalf("could not get text content: %v", err)
		}
		fmt.Println(textContent)
	}
}

func getUrlCoordinates(mapUrl string) pgtype.Point {
	u, err := url.Parse(mapUrl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	queryParams := u.Query()
	coordString := queryParams.Get("q")
	coordArr := strings.Split(coordString, ",")
	var point pgtype.Point
	x, err := strconv.ParseFloat(coordArr[0], 64)
	if err != nil {
		log.Fatalf("could not parse float to point.P.X: %v", err)
	}
	y, err := strconv.ParseFloat(coordArr[1], 64)
	if err != nil {
		log.Fatalf("could not parse float to point.P.X: %v", err)
	}
	point = pgtype.Point{
		P:     pgtype.Vec2{X: x, Y: y},
		Valid: true,
	}
	return point
}

func takeScreenShot(page playwright.Page, path string, fullPage bool) {
	// Take screenshot of full page
	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String("path"),
		FullPage: playwright.Bool(fullPage),
	})
	if err != nil {
		log.Fatalf("could not take screenshot: %v", err)
	}
}

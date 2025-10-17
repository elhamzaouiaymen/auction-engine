package main

import (
	"log"

	"github.com/playwright-community/playwright-go"
)

func LogError(errorStr string, err error) {
	log.Fatalln(errorStr, err)
}

func main() {

	pw, err := playwright.Run()

	if err != nil {
		LogError("could not start playwright", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		LogError("could not launch browser", err)
	}

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		BaseURL: playwright.String("https://www.troostwijkauctions.com/"),
	})

	if err != nil {
		LogError("could not create page", err)
	}
	_, err = page.Goto("https://www.troostwijkauctions.com/")
	if err != nil {
		LogError("could not goto page", err)
	}

	pageTitle, err := page.Title()
	if err != nil {
		LogError("could not get page title", err)
	}
	log.Println("Page title:", pageTitle)

	// search for "construction" and get all results
	search := page.Locator("[data-cy='header-search-input']")

	err = search.Fill("iphone")
	if err != nil {
		LogError("could not fill search box", err)
	}

	err = search.Press("Enter")
	if err != nil {
		LogError("could not press enter", err)
	}

	results := make([]string, 0)
	// paginate through results
	// as long as the next button 
	// is visible, click it and 
	// get the results
	for {
		hasNext, err := page.Locator("[data-cy='pagination-next-link']").IsVisible()
		if err != nil {
			LogError("could not check for next button", err)
		}
		if !hasNext {
			break
		}

		// get all results on page
		items, err := page.Locator("[data-cy='lot-card']").All()
		if err != nil {
			LogError("could not get lot cards", err)
		}

		for _, item := range items {
			resultText, err := item.TextContent()
			if err != nil {
				LogError("could not get lot card text", err)
			}
			results = append(results, resultText)
		}

		next := page.Locator("a[data-cy='pagination-next-link']").First()
		if err != nil {
			LogError("could not get next button", err)
		}
		err = next.ScrollIntoViewIfNeeded()

		if err != nil {
			LogError("could not click next button", err)
		}

		link, err := next.GetAttribute("href")
		if err != nil {
			LogError("could not get next button href", err)
		}

		log.Println("Navigating to next page:", link)

		_, err = page.Goto(link)

		// err = next.Click()
		// if err != nil {
		// 	LogError("could not click next button", err)
		// }

		// // wait for page to load
		// time.Sleep(2 * time.Second)
	} 

}
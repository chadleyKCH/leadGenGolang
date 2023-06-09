package search

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
)

var (
	StateAbb  string
	DriverURL string
	Lead      string
)

func SearchThomasnet() {

	const (
		// These paths will be different on your system.
		// chromeDriverPath = "C:/Users/coleh/VS_CODES/LG_Other/chromedriver.exe"
		chromeDriverPath = "/usr/lib/chromium/chromedriver"

		// chromeDriverPath = "/usr/local/bin/chromedriver"
		port = 4445
	)

	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
		// selenium.Output(os.Stderr),
	}
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{
		"browserName": "chrome",
		"goog:chromeOptions": map[string]interface{}{
			"args": []string{
				"--headless",
				"--disable-gpu",
				"--no-sandbox",
				"--disable-dev-shm-usage",
			},
		},
	}
	P := 4445
	driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", P))
	if err != nil {
		log.Fatal("COULDN'T START NEW SELENIUM REMOTE")
	}
	driver.Get("https://www.thomasnet.com/")
	time.Sleep(time.Second * 8)
	// Cookies
	c, _ := driver.FindElement(selenium.ByCSSSelector, "#gdpr-btn-accept")
	if c != nil {
		c.Click()
	} else {
		fmt.Println("NO COOKIES")
	}
	// Select the Search bar
	input, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > div > div.site-search__search-query-input-wrap.search-suggest-preview > input")
	if err != nil {
		log.Fatalf("FAILED TO FIND THE SEARCH BAR: %v", err)
	}
	//Enter lead from Lead_Template.xlsx into the search bar
	err = input.SendKeys(Lead)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
	// Define the button selector and search for its element.
	search, err := driver.FindElement(selenium.ByCSSSelector, "#homesearch > form > div > button")
	if err != nil {
		panic(err)
	}
	//Click on the search button.
	err = search.Click()
	if err != nil {
		log.Fatalf("FAILED TO CLICK THE SEARCH BUTTON: %v", err)
	}
	time.Sleep(time.Second * 3)
	// Looks for iframe
	// iframeElem, err := driver.FindElement(selenium.ByCSSSelector, "iframe[src*='about:blank']")
	// if err != nil {
	// 	log.Fatalf("FAILED TO FIND THE IFRAME ELEMENT: %v", err)
	// }
	// // Switches to iframe
	// err = driver.SwitchFrame(iframeElem)
	// if err != nil {
	// 	log.Fatalf("FAILED TO SWITCH TO THE IFRAME: %v", err)
	// }
	// time.Sleep(time.Second * 2)
	// // Looks for accept button
	// acceptNoti, err := driver.FindElement(selenium.ByXPATH, "/html/body/appcues/cue/section/div/div[3]/div/div/div/div/div/a")
	// if err != nil {
	// 	log.Fatalf("FAILED TO FIND THE IFRAME ACCEPT BUTTON: %v", err)
	// }
	// // Clicks accept button
	// err = acceptNoti.Click()
	// if err != nil {
	// 	log.Fatalf("FAILED TO CLICK THE IFRAME ACCEPT BUTTON: %v", err)
	// }
	// // Switches back to main frame
	// switchDefault := driver.SwitchFrame(nil)
	// if switchDefault != nil {
	// 	log.Fatalf("FAILED TO SWITCH BACK TO MAIN FRAME: %v", switchDefault)
	// }
	// time.Sleep(time.Second * 2)

	// Checks if StateAbb is blank. If so, then return the driver URL to scrape.go
	if StateAbb == "" {
		DriverURL, _ = driver.CurrentURL()
		return
	} else {
		// Find the region dropdown
		regionDropdown, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > a")
		if err != nil {
			log.Fatalf("FAILED TO FIND SELECT REGION DROPDOWN: %v", err)
		}
		// Click the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("FAILED TO CLICK SELECT REGION DROPDOWN: %v", err)
		}
		time.Sleep(time.Second * 3)

		// Finds specified region
		regionSelect, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > div > div.thm-custom-select.search-options-regions > div [data-value="+StateAbb+"]")
		if err != nil {
			log.Fatalf("COULD NOT FIND THE stateabb SPECIFIED: %v", err)
		}
		// Clicks specified region
		err = regionSelect.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK THE stateabb: %v", err)
		}
		// Select the region dropdown
		err = regionDropdown.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK THE REGION DROPDOWN: %v", err)
		}
		time.Sleep(time.Second * 3)
		// Finds search button
		regionSearch, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.logged-out > header > div.site-header__section > div > div.site-header__section-header__utility > form > div > button")
		if err != nil {
			log.Fatalf("COULD NOT FIND SEARCH BUTTON AFTER THE REGION HAS BEEN SELECTED: %v", err)
		}
		// Clicks search
		err = regionSearch.Click()
		if err != nil {
			log.Fatalf("COULD NOT CLICK SEARCH BUTTON AFTER THE REGION HAS BEEN SELECTED: %v", err)
		}
		time.Sleep(time.Second * 3)

		results, err := driver.FindElement(selenium.ByCSSSelector, "body > div.site-wrap.interim-search-results.logged-out > section.network-search-results > div > div.network-search-results__primary > div > section:nth-child(1) > div > table > tbody > tr:nth-child(1) > td > a")
		if err != nil {
			fmt.Printf("COULD NOT FIND NETWORK RESULT: %v "+"MOVING ON...", err)
		}
		if results != nil {
			err = results.Click()
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Printf("ELEMENT NETWORK NOT FOUND: %v "+"MOVING ON...", err)
		}
		// Finds the 'Located in' option
		LocIn, err := driver.FindElement(selenium.ByCSSSelector, "#main > div.filter-block.located-serving-card > ul > li:nth-child(1) > a")

		time.Sleep(time.Second * 3)
		// If the 'Located in' option does not exist then it skips over...else: clicks located in
		if LocIn != nil {
			// Clicks 'Located in'
			err = LocIn.Click()
			if err != nil {
				log.Fatalf("COULD NOT CLICK ON 'LOCATED IN': %v", err)
			}
		} else {
			log.Printf("COULD NOT FIND 'LOCATED IN' ELEMENT: %v", err)
		}
		time.Sleep(time.Second * 3)

		DriverURL, _ = driver.CurrentURL()
	}

}

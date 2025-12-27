package bot

import (
	"fmt"
	"time"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/input"
	"linkedin-bot/internal/human"
)

type LinkedInBot struct {
	Browser *rod.Browser
	Page    *rod.Page
	Logger  func(string)
}

func NewBot(logger func(string)) *LinkedInBot {
	logger("Launching Stealth Browser...")
	path, _ := launcher.LookPath()
	logger("Chrome path: " + path)
	u := launcher.New().
		Bin(path).
		Headless(false).
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-sandbox").
		Set("no-sandbox").
		Set("disable-dev-shm-usage").
		Set("disable-gpu").
		MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	
	time.Sleep(2 * time.Second)
	
	logger("Creating page...")
	page := browser.MustPage("about:blank")
	logger("Page created successfully")
	
	err := page.Navigate("https://www.linkedin.com")
	if err != nil {
		logger(fmt.Sprintf("ERROR: %v", err))
	}
	page.MustWaitLoad()

	return &LinkedInBot{Browser: browser, Page: page, Logger: logger}
}

func (b *LinkedInBot) Login(username, password string) {
    b.Logger("Navigating to LinkedIn Login...")
    err := b.Page.Navigate("https://www.linkedin.com/login")
    if err != nil {
        b.Logger(fmt.Sprintf("ERROR navigating to login: %v", err))
        return
    }
    b.Page.MustWaitLoad()
    time.Sleep(2 * time.Second)

    b.Logger("Typing Credentials (Stealth Mode)...")
    
    // Try different selectors for username field
    usernameEl, err := b.Page.Element("input[name='session_key']")
    if err != nil {
        b.Logger("ERROR: Could not find username field")
        return
    }
    human.TypeStealthily(usernameEl, username)
    time.Sleep(500 * time.Millisecond)
    b.Page.KeyActions().Press(input.Tab).MustDo()
    
    passwordEl, err := b.Page.Element("input[name='session_password']")
    if err != nil {
        b.Logger("ERROR: Could not find password field")
        return
    }
    human.TypeStealthily(passwordEl, password)
    time.Sleep(500 * time.Millisecond)

    submitBtn, err := b.Page.Element("button[type='submit']")
    if err != nil {
        b.Logger("ERROR: Could not find submit button")
        return
    }
    human.MoveMouseNaturally(b.Page, submitBtn)
    submitBtn.MustClick()
    
    b.Logger("Login Submitted. Waiting...")
    b.Page.MustWaitLoad()
}

// *** NEW FUNCTION ADDED HERE ***
func (b *LinkedInBot) SearchAndConnect(jobTitle string) {
    b.Logger(fmt.Sprintf("Searching for: %s", jobTitle))
    
    // 1. Navigate to Search Results
    searchURL := fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s", jobTitle)
    err := b.Page.Navigate(searchURL)
    if err != nil {
        b.Logger(fmt.Sprintf("ERROR navigating to search: %v", err))
        return
    }
    b.Page.MustWaitStable()

    b.Logger("Scanning for candidates...")

    // 2. Find all 'Connect' buttons
    buttons, err := b.Page.Elements("button")
    if err != nil {
        b.Logger(fmt.Sprintf("ERROR finding buttons: %v", err))
        return
    }
    
    dailyLimit := 5
    sentCount := 0

    for _, btn := range buttons {
        if sentCount >= dailyLimit {
            b.Logger("Daily limit reached. Stopping.")
            break
        }

        if text, _ := btn.Text(); text == "Connect" {
            btn.MustScrollIntoView()
            
            b.Logger("Clicking Connect...")
            human.MoveMouseNaturally(b.Page, btn)
            btn.MustClick()
            time.Sleep(1 * time.Second)

            // 3. Handle "Add a note" popup
            addNoteBtn, err := b.Page.ElementR("button", "Add a note")
            if err == nil {
                b.Logger("Adding personalized note...")
                addNoteBtn.MustClick()
                time.Sleep(1 * time.Second)

                // Type message
                textArea, err := b.Page.Element("textarea[name='message']")
                if err == nil {
                    message := fmt.Sprintf("Hi, I noticed your work in %s and would love to connect!", jobTitle)
                    human.TypeStealthily(textArea, message)
                    time.Sleep(1 * time.Second)

                    // Click Send
                    sendBtn, err := b.Page.ElementR("button", "Send")
                    if err == nil {
                        human.MoveMouseNaturally(b.Page, sendBtn)
                        sendBtn.MustClick()
                        sentCount++
                        b.Logger(fmt.Sprintf("Invite sent! (%d/%d)", sentCount, dailyLimit))
                    }
                }
            } else {
                // Fallback if no note option
                if sendBtn, err := b.Page.ElementR("button", "Send"); err == nil {
                    sendBtn.MustClick()
                    sentCount++
                    b.Logger("Sent without note.")
                }
            }

            b.Logger("Cooling down (5s)...")
            time.Sleep(5 * time.Second)
        }
    }
    
    if sentCount == 0 {
        b.Logger("No connectable profiles found on this page.")
    } else {
        b.Logger("Batch complete.")
    }
}
func (b *LinkedInBot) Stop() {
	if b.Browser != nil { b.Browser.MustClose() }
}
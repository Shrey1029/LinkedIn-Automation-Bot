package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"linkedin-bot/internal/bot"
)

var (
	logs      []string
	logsMutex sync.Mutex
	activeBot *bot.LinkedInBot
)

func AddLog(msg string) {
	logsMutex.Lock()
	defer logsMutex.Unlock()
	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf("[%s] %s", timestamp, msg)
	logs = append(logs, formatted)
	if len(logs) > 100 {
		logs = logs[len(logs)-100:]
	}
	fmt.Println(formatted)
}

func main() {
	godotenv.Load()
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.File("web/templates/index.html")
	})

	e.POST("/start", func(c echo.Context) error {
		query := c.FormValue("query")
		AddLog(fmt.Sprintf("Received Start Command. Target: %s", query))
		go func() {
			user := os.Getenv("LINKEDIN_USER")
			pass := os.Getenv("LINKEDIN_PASS")
			if user == "" || pass == "" {
				AddLog("ERROR: Credentials missing in .env file")
				return
			}
			activeBot = bot.NewBot(AddLog)
			if activeBot == nil {
				AddLog("ERROR: Failed to create bot")
				return
			}
			
			// 1. Login
			activeBot.Login(user, pass)
			
			// 2. *** START SEARCHING (This was missing) ***
			if activeBot != nil {
				activeBot.SearchAndConnect(query)
			}
			
			AddLog("Workflow Complete.")
		}()
		return c.String(http.StatusOK, "Started")
	})

	e.POST("/stop", func(c echo.Context) error {
		AddLog("Stopping Bot...")
		if activeBot != nil {
			activeBot.Stop()
			activeBot = nil
		}
		return c.String(http.StatusOK, "Stopped")
	})

	e.GET("/logs", func(c echo.Context) error {
		logsMutex.Lock()
		defer logsMutex.Unlock()
		html := ""
		for i := len(logs) - 1; i >= 0; i-- {
			html += fmt.Sprintf("<div class='text-green-400'>%s</div>", logs[i])
		}
		return c.HTML(http.StatusOK, html)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
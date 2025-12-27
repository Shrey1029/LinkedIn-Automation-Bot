# LinkedIn Stealth Automation Suite

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Modular-green)
![Status](https://img.shields.io/badge/Status-Proof%20of%20Concept-orange)

A sophisticated, modular automation tool for LinkedIn built with **Go** and **Rod**. This project demonstrates advanced browser automation techniques, featuring a **human-behavior simulation engine** to bypass bot detection filters. It includes a real-time web dashboard powered by **HTMX**.

---

## ⚠️ Legal Disclaimer

**EDUCATIONAL PURPOSE ONLY.**
This software is a Proof of Concept (PoC) designed to demonstrate technical capabilities in browser automation and stealth engineering. Automating LinkedIn violates their [User Agreement](https://www.linkedin.com/legal/user-agreement). The author is not responsible for account bans or legal consequences resulting from the use of this tool. **Do not use on primary accounts.**

---

## 🌟 Key Features

### 🤖 Automation Engine
* **Headless-Optional Browser:** Runs a controlled instance of Chrome/Chromium.
* **Smart Login:** Handles authentication and session stability.
* **Targeted Outreach:** Searches for specific job titles (e.g., "Software Engineer") and sends connection requests.
* **Follow-Up System:** Scans recent connections and sends personalized follow-up messages using dynamic templates.

### 🥷 Advanced Stealth (Anti-Detection)
* **Bézier Curve Mouse Movement:** Instead of linear robotic paths, the mouse cursor follows randomized cubic Bézier curves with variable speed and micro-jitters, simulating human hand motor control.
* **Human-Like Typing:** Keystrokes are injected with randomized delays (50ms–150ms) to mimic natural typing rhythm.
* **Browser Fingerprint Masking:** Uses `rod-stealth` to strip `navigator.webdriver` flags and randomize viewport parameters.
* **Rate Limiting:** Enforces strict daily caps (e.g., 5 requests/day) to prevent flagging.

### 🖥️ Command Center UI
* **Dark Mode Dashboard:** A professional interface built with Tailwind CSS.
* **Real-Time Logging:** Uses HTMX polling to stream server-side logs to the browser instantly.
* **No Node.js Required:** The entire frontend is served directly from the Go binary.

---

### Core Components
* The Brain (cmd/server): Orchestrates the concurrent Goroutines. It spins up the HTTP server and manages the communication between the UI and the Bot Engine via a thread-safe logging channel.

* The Bot (internal/bot): Manages the browser lifecycle. It tracks state, handles DOM queries, and manages the JSON-based history file to prevent duplicate messaging.

* The Hand (internal/human): A pure math package. It calculates trajectory points for mouse movement and handles "sleep" randomization to defeat heuristic analysis.

## 🏗️ Technical Architecture

The project follows a modular **Clean Architecture** pattern to ensure maintainability and separation of concerns.

```text
linkedin-bot/
├── cmd/
│   └── server/
│       └── main.go           # Application Entry Point & HTTP Server
├── internal/
│   ├── bot/
│   │   └── engine.go         # Core Automation Logic (Login, Search, Message)
│   └── human/
│       └── mouse.go          # Stealth Physics Engine (Bézier Curves)
├── web/
│   └── templates/
│       └── index.html        # HTMX + Tailwind Dashboard
├── history.json              # Local persistence for tracking sent messages
├── .env                      # Configuration credentials
├── go.mod                    # Dependency definitions
└── README.md                 # Project Documentation 

# SilentFetch
Exploits live-tracker site via backend bypass 😉


# 📱 Phone Info Lookup Tool (Golang)

**By Abdul Ahad – JUST HACKED ON 👽**

A powerful CLI tool that scrapes phone number details from a tracker website. Bypasses frontend restrictions and scrapes unlimited data — made for ethical hacking learning and research purposes.

---

## 🛠️ Why Use This Tool Instead of the Website?

The original website (`https://live-tracker.site`) is limited in functionality:
- 🔒 Accepts only **one number at a time**
- ❌ No option for **bulk lookup**
- 🧱 Frontend restriction to slow down users
- 🕵️ No CAPTCHA or bot protection

This tool **bypasses all limitations** by directly sending crafted POST requests to the backend and parsing the result.

---

## 🔓 What Vulnerability Does It Exploit?

This tool takes advantage of:

### 1. **Improper Rate Limiting**
- The server does not enforce request limits.
- Allows sending hundreds of queries in a short time.
- No throttling or IP-based restrictions.

### 2. **Lack of Bot Detection / CAPTCHA**
- No CAPTCHA or JavaScript challenge (e.g., Cloudflare).
- Can be easily scripted and automated.

### 3. **Potential Sensitive Data Disclosure**
- Information like name, address, CNIC is publicly exposed without authentication.

> ⚠️ **This tool is for educational purposes only. Do not use on live systems without permission.**

---

## 🚀 Features

- 🔍 Search a single number or a full list
- 🎨 Formatted colored output
- 📁 Auto-save results to `results.txt`
- 🚀 Fast and reliable scraper using Go's concurrency
- 🧠 ASCII hacker vibe banner at start

---

## 📸 Sample Output & Usage

```bash
🌟 Entry #1 | Just Hacked On 👽
╔════════════════╦══════════════════════════════════════════════════════════╗
║ 🔖 Field       ║ 📝 Value                                                 ║
╠════════════════╬══════════════════════════════════════════════════════════╣
║ Name           ║ JUST-HACKED-ON                                                ║
║ CNIC           ║ 35201-1234567-8                                         ║
║ Address        ║ Lahore, Punjab                                          ║
╚════════════════╩══════════════════════════════════════════════════════════╝

⚙️ Installation
1. Clone this repo

git clone https://github.com/justhackedon/TrackBreaker.git
cd TrackBreaker

2. Install dependencies

go get github.com/PuerkitoBio/goquery
go get github.com/fatih/color
go get github.com/schollz/progressbar/v3

3. Run the tool

go run main.go -num 03001234567

go run main.go -l numbers.txt

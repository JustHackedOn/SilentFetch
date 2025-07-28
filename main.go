package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// ASCII banner for JUST HACKED ON
const banner = `
     ██╗██╗   ██╗███████╗████████╗    ██╗  ██╗ █████╗  ██████╗██╗  ██╗███████╗██████╗      ██████╗ ███╗   ██╗
     ██║██║   ██║██╔════╝╚══██╔══╝    ██║  ██║██╔══██╗██╔════╝██║ ██╔╝██╔════╝██╔══██╗    ██╔═══██╗████╗  ██║
     ██║██║   ██║███████╗   ██║       ███████║███████║██║     █████╔╝ █████╗  ██║  ██║    ██║   ██║██╔██╗ ██║
██   ██║██║   ██║╚════██║   ██║       ██╔══██║██╔══██║██║     ██╔═██╗ ██╔══╝  ██║  ██║    ██║   ██║██║╚██╗██║
╚█████╔╝╚██████╔╝███████║   ██║       ██║  ██║██║  ██║╚██████╗██║  ██╗███████╗██████╔╝    ╚██████╔╝██║ ╚████║
 ╚════╝  ╚═════╝ ╚══════╝   ╚═╝       ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═════╝      ╚═════╝ ╚═╝  ╚═══╝
                                                                   Abdul Ahad   ==> Security Just an illusion
`

// Result represents a single key-value pair from the scraped data
type Result struct {
	Field string
	Value string
}

// wrapText wraps text to a specified width
func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Split(text, " ")
	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			if len(word) > width {
				for len(word) > width {
					lines = append(lines, word[:width])
					word = word[width:]
				}
			}
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// FetchInfo scrapes data for a given phone number
func fetchInfo(number string) ([]Result, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	data := url.Values{}
	data.Set("searchinfo", number)

	req, err := http.NewRequest("POST", "https://live-tracker.site/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, nil // Silent error
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://live-tracker.site")
	req.Header.Set("Referer", "https://live-tracker.site/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil // Silent error
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil // Silent error
	}

	var results []Result
	doc.Find("div.resultcontainer").Each(func(i int, container *goquery.Selection) {
		container.Find("div.row").Each(func(j int, row *goquery.Selection) {
			key := strings.TrimSuffix(row.Find("span.detailshead").Text(), " :")
			value := strings.TrimSpace(row.Find("span.details").Text())
			if key != "" && value != "" {
				results = append(results, Result{Field: key, Value: value})
			}
		})
	})

	return results, nil
}

// PrintResult displays results in a formatted table
func printResult(number string, results []Result) {
	if len(results) == 0 {
		color.Red("❌ No records found for %s", number)
		return
	}

	color.Green("\n📱 Results for: %s", number)
	entryCount := 0
	for i, result := range results {
		if i == 0 || results[i-1].Field == "Address" || (i > 0 && results[i-1].Field == "CNIC" && result.Field == "Name") {
			if i > 0 {
				color.Yellow("╚════════════════╩══════════════════════════════════════════════════════════╝")
			}
			entryCount++
			color.Magenta("\n🌟 Entry #%d | Just Hacked On 👽\n", entryCount)
			color.Yellow("╔════════════════╦══════════════════════════════════════════════════════════╗")
			color.Cyan("║ 🔖 Field       ║ 📝 Value                                                 ║")
			color.Yellow("╠════════════════╬══════════════════════════════════════════════════════════╣")
		}
		color.White("║ %-14s ║ %-60s ║", result.Field, result.Value)
		if i == len(results)-1 || (i < len(results)-1 && results[i+1].Field == "Name" && results[i].Field == "Address") {
			color.Yellow("╚════════════════╩══════════════════════════════════════════════════════════╝")
		}
	}
}

// SaveToFile writes results to a file
func saveToFile(number string, results []Result, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil // Silent error
	}
	defer f.Close()

	if len(results) == 0 {
		_, err := f.WriteString(fmt.Sprintf("No records found for %s\n\n", number))
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("\nResults for %s\n", number))
	if err != nil {
		return nil // Silent error
	}

	entryCount := 0
	for i, result := range results {
		if i == 0 || (i > 0 && result.Field == "Name" && (results[i-1].Field == "Address" || results[i-1].Field == "CNIC")) {
			entryCount++
			_, err = f.WriteString(fmt.Sprintf("--- Entry #%d ---\n", entryCount))
			if err != nil {
				return nil // Silent error
			}
		}
		_, err = f.WriteString(fmt.Sprintf("%s: %s\n", result.Field, result.Value))
		if err != nil {
			return nil // Silent error
		}
	}
	_, err = f.WriteString("\n")
	return err
}

func main() {
	// Print ASCII banner
	color.Cyan(banner)

	// CLI arguments
	numPtr := flag.String("num", "", "Phone number to search")
	listPtr := flag.String("l", "", "File path to a list of numbers")
	flag.Parse()

	if *numPtr == "" && *listPtr == "" {
		color.Red("❗ Please provide -num OR -l argument")
		return
	}

	// Mode 1: Single Number
	if *numPtr != "" {
		results, _ := fetchInfo(*numPtr)
		printResult(*numPtr, results)
		return
	}

	// Mode 2: List of Numbers
	if _, err := os.Stat(*listPtr); os.IsNotExist(err) {
		color.Red("❌ File not found: %s", *listPtr)
		return
	}

	outputFile := "results.txt"
	os.Remove(outputFile) // Remove old output file

	file, err := os.Open(*listPtr)
	if err != nil {
		color.Red("❌ Error opening file: %s", *listPtr)
		return
	}
	defer file.Close()

	var numbers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			numbers = append(numbers, line)
		}
	}

	bar := progressbar.NewOptions(len(numbers),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionSetDescription("📡 Scanning..."),
		progressbar.OptionSetWidth(30),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetRenderBlankState(true),
	)

	for _, number := range numbers {
		results, _ := fetchInfo(number)
		if err := saveToFile(number, results, outputFile); err != nil {
			continue // Silent error
		}
		bar.Add(1)
	}

	color.Green("\n✅ Saved all results to %s", outputFile)
}

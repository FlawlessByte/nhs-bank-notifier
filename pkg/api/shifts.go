// pkg/api/shifts.go
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nhs-bank-notifier/pkg/config"
	"nhs-bank-notifier/pkg/logger"
	"nhs-bank-notifier/pkg/notifier"
	"strings"
	"time"

	"github.com/jftuga/TtlMap"
)

type Shift struct {
	Unit      string `json:"unit"`
	Shift     string `json:"shift"`
	Date      string `json:"onDate"`
	Id        int    `json:"dutyId"`
	RequestId string `json:"requestId"`
}

type ResponseData struct {
	TotalItemsCount int     `json:"TotalItemsCount"`
	VacantDuties    []Shift `json:"vacantDuties"`
}

func Login(httpClient *http.Client, loginURL, username, password string) error {

	log := logger.GetLogger()

	loginPayload := url.Values{
		"ShowResetPasswordLink": {"True"},
		"ReturnUrl":             {""},
		"ApplicationName":       {""},
		"Username":              {username},
		"Password":              {password},
		"btnLogin":              {"Log in"},
	}

	req, err := http.NewRequest("POST", loginURL, bytes.NewBufferString(loginPayload.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %s", resp.Status)
	}

	log.Info("Logged in successfully!")
	return nil
}

func FetchShifts(httpClient *http.Client, ttlMap *TtlMap.TtlMap[string], cfg *config.Config) []Shift {
	log := logger.GetLogger()

	// Calculate date range
	now := time.Now()
	startDate := now.Format("2006/01/02")
	endDate := now.AddDate(0, 0, 7).Format("2006/01/02")
	log.Debug("Fetching shifts from ", startDate, " to ", endDate)

	// Step 2: Make API request to fetch shifts
	apiURL := "https://ich.allocate-cloud.co.uk/EmployeeOnlineHealth/ICHLIVE/Roster/BankShifts/UnfilledShifts"
	payload := map[string]interface{}{
		"take":     15,
		"skip":     0,
		"page":     1,
		"pageSize": 100,
		"from":     startDate,
		"to":       endDate,
	}
	log.Debug("Fetching shifts with payload: ", payload)

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Errorf("Failed to create request: %v", err)
	}

	// Set headers to mimic a real browser request
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-IN,en-GB;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://ich.allocate-cloud.co.uk")
	req.Header.Set("Referer", "https://ich.allocate-cloud.co.uk/EmployeeOnlineHealth/ICHLIVE/Roster/BankShifts")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	req.Header.Set("Sec-Ch-Ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch shifts: %s", resp.Status)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		log.Warn("Session expired, re-logging in...")
		if loginErr := Login(httpClient, cfg.LoginURL, cfg.NHSUsername, cfg.NHSPassword); loginErr != nil {
			log.Errorf("Re-login failed: %v", loginErr)
			if err := notifier.SendTelegramMessage(cfg.TelegramToken, cfg.TelegramChatID, "Re-login failed. Inform the creator!"); err != nil {
				log.Errorf("Error sending Telegram message: %v", err)
			}
			panic(loginErr)
		}
		return []Shift{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseData ResponseData
	if err := json.Unmarshal(body, &responseData); err != nil {
		log.Fatal(err)
	}

	log.Debugf("Response Status: %s", resp.Status)

	// Filter shifts from a particular units
	var shifts []Shift
	for _, shift := range responseData.VacantDuties {
		if strings.HasPrefix(shift.Unit, cfg.NHSUnit) {
			if ttlMap.Get(fmt.Sprintf("%d", shift.Id)) == nil {
				shifts = append(shifts, shift)
				ttlMap.Put(fmt.Sprintf("%d", shift.Id), shift.Id)
			} else {
				log.Debug("Shift already notified & skipping: ", shift.Id)
			}
		}
	}

	return shifts
}

func FormatShiftsMessage(shifts []Shift) string {
	message := "New shifts available:\n"
	for _, shift := range shifts {
		message += fmt.Sprintf("- Unit: %s, Shift: %s, Date: %s\n", shift.Unit, shift.Shift, shift.Date) // Customize this message as needed
	}
	return message
}

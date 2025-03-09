package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	clmidLogin  = "CLMAuthLoginRequest"
	clmidLogout = "CLMAuthLogoutRequest"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags) // 標準出力にログ出力、タイムスタンプ付き
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file not found", err)
	}

	// デモ環境
	// demoBaseURL := os.Getenv("TACHIBANA_DEMO_URL")
	// demoUserID := os.Getenv("TACHIBANA_USER_ID")
	// demoPassword := os.Getenv("TACHIBANA_PASSWORD")
	// runEnvironment("Demo", demoBaseURL, demoUserID, demoPassword)

	// fmt.Println("------------------------------------") // 環境の区切り

	// 本番環境
	prodBaseURL := os.Getenv("TACHIBANA_BASE_URL")
	prodUserID := os.Getenv("TACHIBANA_USER_ID")
	prodPassword := os.Getenv("TACHIBANA_PASSWORD")
	runEnvironment("Production", prodBaseURL, prodUserID, prodPassword)
}

func runEnvironment(envName, baseURL, userID, password string) {
	logger.Printf("===== %s Environment =====", envName)
	if baseURL == "" || userID == "" || password == "" {
		logger.Fatalf("Error: %s environment variables are not set.", envName)
		return // Fatalf の場合は、ここで処理を終了
	}

	requestURL, pNo, err := login(baseURL, userID, password)
	if err != nil {
		logger.Printf("%s Login Error: %v\n", envName, err)
		return // エラーの場合は、ログアウト処理に進まない
	}

	logger.Printf("---------- %s Logout ----------", envName)
	if err := logout(requestURL, pNo); err != nil {
		logger.Printf("%s Logout Error: %v\n", envName, err)
	}
}

func login(baseURL, userID, password string) (string, int64, error) {
	authURL, _ := url.JoinPath(baseURL, "auth/")
	payload := map[string]string{
		"p_no":      "1",
		"p_sd_date": time.Now().Format("2006.01.02-15:04:05.000"),
		"sPassword": password,
		"sUserId":   userID,
		"sCLMID":    clmidLogin,
		"sJsonOfmt": "4",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("JSON encode error: %w", err)
	}

	encodedPayload := url.QueryEscape(string(payloadJSON))
	authURL += "?" + encodedPayload

	logger.Printf("Login Request URL: %s", authURL)
	logger.Printf("Login Request Payload: %s", string(payloadJSON))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, authURL, nil)
	if err != nil {
		return "", 0, fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("request send error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("response body read error: %w", err)
	}

	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
	if err != nil {
		return "", 0, fmt.Errorf("shift-jis decode error: %w", err)
	}

	logger.Printf("Status Code: %d", resp.StatusCode)

	// JSON レスポンスを整形して出力
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyUTF8, "", "    "); err != nil {
		// 整形に失敗しても、元のレスポンスは表示
		logger.Printf("Response Body (Unformatted):\n%s", string(bodyUTF8))
	} else {
		logger.Printf("Response Body (Formatted):\n%s", prettyJSON.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyUTF8, &response); err != nil {
		return "", 0, fmt.Errorf("JSON decode error: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); !ok {
		// sResultCode が存在しない場合
		return "", 0, fmt.Errorf("API error: sResultCode not found in response")
	} else if resultCode != "0" {
		// sResultCode が "0" 以外の場合 (エラー)
		sResultText := ""
		if rt, ok := response["sResultText"].(string); ok {
			sResultText = rt
		}

		if warnCode, ok := response["sWarningCode"].(string); ok {
			sWarningText := ""
			if wt, ok := response["sWarningText"].(string); ok {
				sWarningText = wt
			}
			return "", 0, fmt.Errorf("API error: ResultCode=%s, ResultText=%s, WarningCode=%s, WarningText=%s",
				resultCode, sResultText, warnCode, sWarningText)
		}

		return "", 0, fmt.Errorf("API error: ResultCode=%s, ResultText=%s", resultCode, sResultText)
	}

	requestURL, ok := response["sUrlRequest"].(string)
	if !ok {
		return "", 0, fmt.Errorf("sUrlRequest not found in response")
	}

	pNoStr, ok := response["p_no"].(string)
	if !ok {
		return "", 0, fmt.Errorf("p_no not found in response")
	}
	pNo, err := strconv.ParseInt(pNoStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("failed to parse p_no: %w", err)
	}
	return requestURL, pNo, nil
}

func logout(requestURL string, pNo int64) error {
	pNo++

	payload := map[string]string{
		"sCLMID":    clmidLogout,
		"p_no":      strconv.FormatInt(pNo, 10),
		"p_sd_date": time.Now().Format("2006.01.02-15:04:05.000"),
		"sJsonOfmt": "4",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON encode error: %w", err)
	}

	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL += "?" + encodedPayload

	logger.Printf("Logout Request URL: %s", requestURL)
	logger.Printf("Logout Request Payload: %s", string(payloadJSON))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request send error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("response body read error: %w", err)
	}

	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
	if err != nil {
		return fmt.Errorf("shift-jis decode error: %w", err)
	}

	logger.Printf("Status Code: %d", resp.StatusCode)

	// JSON レスポンスを整形して出力
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bodyUTF8, "", "    "); err != nil {
		logger.Printf("Response Body (Unformatted):\n%s", string(bodyUTF8))
	} else {
		logger.Printf("Response Body (Formatted):\n%s", prettyJSON.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyUTF8, &response); err != nil {
		return fmt.Errorf("JSON decode error: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); !ok {
		// sResultCode が存在しない場合
		return fmt.Errorf("API error: sResultCode not found in response")
	} else if resultCode != "0" {
		sResultText := ""
		if rt, ok := response["sResultText"].(string); ok {
			sResultText = rt
		}
		return fmt.Errorf("API error: ResultCode=%s, ResultText=%s", resultCode, sResultText)
	}

	return nil
}

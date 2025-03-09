package main

import (
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

func main() {
	// .env ファイルから環境変数を読み込む
	err := godotenv.Load(".env") // .envファイルはプロジェクトルートに配置
	if err != nil {
		log.Fatal(".env file not found", err) //log.Fatalへ変更
	}

	// デモ環境でのログインとログアウト
	fmt.Println("===== Demo Environment =====")
	demoBaseURL := os.Getenv("TACHIBANA_DEMO_URL")
	demoUserID := os.Getenv("TACHIBANA_USER_ID")    // 共通のID
	demoPassword := os.Getenv("TACHIBANA_PASSWORD") // 共通のパスワード
	if demoBaseURL == "" || demoUserID == "" || demoPassword == "" {
		log.Fatal("Error: Demo environment variables are not set.")
	}
	// ログインし、仮想URLと初期p_noを取得
	demoRequestURL, demoPNo, err := login(demoBaseURL, demoUserID, demoPassword)
	if err != nil {
		fmt.Printf("Demo Login Error: %v\n", err)
	} else {
		// 2. ログアウト
		fmt.Println("---------- Demo Logout ----------")
		// デモ環境のログアウト (p_no をインクリメントして使用)
		if err := logout(demoRequestURL, demoPNo); err != nil {
			fmt.Printf("Demo Logout Error: %v\n", err)
		}

	}

	// 本番環境でのログインとログアウト
	fmt.Println("===== Production Environment =====")
	prodBaseURL := os.Getenv("TACHIBANA_BASE_URL")
	prodUserID := os.Getenv("TACHIBANA_USER_ID")    // 共通のID
	prodPassword := os.Getenv("TACHIBANA_PASSWORD") // 共通のパスワード
	if prodBaseURL == "" || prodUserID == "" || prodPassword == "" {
		log.Fatal("Error: Production environment variables are not set.")
	}
	// 本番環境のログイン (p_no は初期化される)
	prodRequestURL, prodPNo, err := login(prodBaseURL, prodUserID, prodPassword)
	if err != nil {
		fmt.Printf("Production Login Error: %v\n", err)
	} else {
		fmt.Println("---------- Prod Logout ----------")
		// 本番環境のログアウト (p_no をインクリメントして使用)
		if err := logout(prodRequestURL, prodPNo); err != nil {
			fmt.Printf("Prod Logout Error: %v\n", err)
		}
	}
}

func login(baseURL, userID, password string) (string, int64, error) {
	// リクエストURLの作成
	authURL, _ := url.JoinPath(baseURL, "auth/") // "auth/" のみ!
	payload := map[string]string{
		"sCLMID":    clmidLogin,
		"sUserId":   userID,
		"sPassword": password,
		"p_sd_date": time.Now().Format("2006.01.02-15:04:05.000"),
		// "p_no":      "1",  //p_noは指定しない
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("JSON encode error: %w", err)
	}

	// URL にクエリパラメータとして JSON全体をエンコードして追加
	authURL += "?" + url.QueryEscape(string(payloadJSON))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, authURL, nil) //GETで作成
	if err != nil {
		return "", 0, fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // 本来は不要なはずだが、念のため残す

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

	// Shift-JIS から UTF-8 に変換
	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
	if err != nil {
		return "", 0, fmt.Errorf("shift-jis decode error: %w", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(bodyUTF8)) // UTF-8で出力

	var response map[string]interface{}
	if err := json.Unmarshal(bodyUTF8, &response); err != nil { //UTF-8でデコード
		return "", 0, fmt.Errorf("JSON decode error: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); !ok || resultCode != "0" {
		if warnCode, ok := response["sWarningCode"].(string); ok {
			return "", 0, fmt.Errorf("API error: %s - %s, Warning: %s - %s", resultCode, response["sResultText"], warnCode, response["sWarningText"])
		}
		return "", 0, fmt.Errorf("API error: %s - %s", resultCode, response["sResultText"])
	}

	// 仮想URL (REQUEST) を取得 (ログアウト時に使用)
	requestURL, ok := response["sUrlRequest"].(string)
	if !ok {
		return "", 0, fmt.Errorf("sUrlRequest not found in response")
	}

	// p_no を取得 (int64に変換)
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

func logout(requestURL string, pNo int64) error { //pNo int64を追加

	//p_noのインクリメント
	pNo++

	payload := map[string]string{
		"sCLMID":    clmidLogout,
		"p_no":      strconv.FormatInt(pNo, 10), //インクリメントしたp_no
		"p_sd_date": time.Now().Format("2006.01.02-15:04:05.000"),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON encode error: %w", err)
	}

	// URL にクエリパラメータとして JSON全体をエンコードして追加
	requestURL += "?" + url.QueryEscape(string(payloadJSON))

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, requestURL, nil) //requestURLではなくbaseURL
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // 本来は不要なはずだが、念のため残す

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

	// Shift-JIS から UTF-8 に変換
	bodyUTF8, _, err := transform.Bytes(japanese.ShiftJIS.NewDecoder(), body)
	if err != nil {
		return fmt.Errorf("shift-jis decode error: %w", err)
	}
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(bodyUTF8)) // UTF-8で出力

	var response map[string]interface{}
	if err := json.Unmarshal(bodyUTF8, &response); err != nil { //UTF-8でデコード
		return fmt.Errorf("JSON decode error: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); !ok || resultCode != "0" {
		return fmt.Errorf("API error: %s - %s", resultCode, response["sResultText"])
	}

	return nil
}

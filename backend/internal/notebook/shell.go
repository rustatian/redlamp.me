package notebook

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

// SessionUser is the signed-in account. The shell embeds it as JSON for the web client.
type SessionUser struct {
	Role     string `json:"role"`
	Username string `json:"username"`
}

const shell = `<!doctype html>
<html lang="ru"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>Тетрадь LeetCode | redlamp</title><link rel="icon" href="/leetcode/assets/icon.svg" type="image/svg+xml"><link rel="stylesheet" href="/leetcode/assets/style.css"><meta name="csp-nonce" nonce="%s">
%s<script type="module" src="/leetcode/assets/main.js"></script></head>
<body><div id="root"></div><noscript>Включите JavaScript, чтобы пользоваться тетрадью.</noscript></body></html>`

// renderShell writes the HTML page that loads the web client. user is nil on the login page.
// The code editor adds <style> elements with the nonce. Each response gets a new random nonce of 128 bits: https://www.w3.org/TR/CSP3/#security-nonces
func renderShell(w http.ResponseWriter, r *http.Request, user *SessionUser) error {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return err
	}
	nonce := base64.StdEncoding.EncodeToString(random)
	session := ""
	if user != nil {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		session = `<script id="session" type="application/json">` + string(data) + `</script>`
	}
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'nonce-"+nonce+"'; connect-src 'self'; media-src 'self' blob:; frame-src 'none'; object-src 'none'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
	send(w, r, http.StatusOK, fmt.Appendf(nil, shell, nonce, session))
	return nil
}

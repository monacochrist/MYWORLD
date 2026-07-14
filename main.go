package main

import (
<<<<<<< HEAD
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
=======
	"crypto/rand"
>>>>>>> origin/master
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
<<<<<<< HEAD
	"strconv"
=======
>>>>>>> origin/master
	"strings"
	"sync"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

type Button struct {
	Text      string  `json:"text"`
	HoverText string  `json:"hover_text"`
	URL       string  `json:"url"`
	Enabled   bool    `json:"enabled"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
}

type AudioTrack struct {
	Name string  `json:"name"`
	URL  string  `json:"url"`
	Loop bool    `json:"loop"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type AlbumTrack struct {
<<<<<<< HEAD
	Name    string `json:"name"`
	URL     string `json:"url"`
	Streams int    `json:"streams"`
=======
	Name string `json:"name"`
	URL  string `json:"url"`
>>>>>>> origin/master
}

type Album struct {
	Name     string       `json:"name"`
	ImageURL string       `json:"image_url"`
	Tracks   []AlbumTrack `json:"tracks"`
	X        float64      `json:"x"`
	Y        float64      `json:"y"`
}

type Config struct {
<<<<<<< HEAD
	Title        string       `json:"title"`
	BgColor      string       `json:"bg_color"`
	TitleColor   string       `json:"title_color"`
	MonthlyUsers int          `json:"monthly_users"`
	TokenBalance int          `json:"token_balance"`
	Buttons      []Button     `json:"buttons"`
	AudioTracks  []AudioTrack `json:"audio_tracks"`
	Albums       []Album      `json:"albums"`
=======
	Title       string       `json:"title"`
	BgColor     string       `json:"bg_color"`
	TitleColor  string       `json:"title_color"`
	Buttons     []Button     `json:"buttons"`
	AudioTracks []AudioTrack `json:"audio_tracks"`
	Albums      []Album      `json:"albums"`
>>>>>>> origin/master
}

type UserData struct {
	Balance         int `json:"balance"`
	SecondsListened int `json:"seconds_listened"`
}

type SSEBroker struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

func NewSSEBroker() *SSEBroker {
	return &SSEBroker{clients: make(map[chan []byte]struct{})}
}

func (b *SSEBroker) Subscribe() chan []byte {
	ch := make(chan []byte, 10)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *SSEBroker) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
}

func (b *SSEBroker) Broadcast(data []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.clients {
		select {
		case ch <- data:
		default:
		}
	}
}

type BalanceStore struct {
	mu       sync.Mutex
	path     string
	balances map[string]UserData
}

func NewBalanceStore(path string) *BalanceStore {
	return &BalanceStore{path: path}
}

func (s *BalanceStore) Lock()   { s.mu.Lock() }
func (s *BalanceStore) Unlock() { s.mu.Unlock() }

func (s *BalanceStore) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.balances = make(map[string]UserData)
			return nil
		}
		return err
	}
	var result map[string]UserData
	if err := json.Unmarshal(data, &result); err == nil {
		s.balances = result
		return nil
	}
	var old map[string]int
	if err := json.Unmarshal(data, &old); err != nil {
		return err
	}
	s.balances = make(map[string]UserData)
	for k, v := range old {
		s.balances[k] = UserData{Balance: v}
	}
	return nil
}

func (s *BalanceStore) Save() error {
	data, err := json.MarshalIndent(s.balances, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *BalanceStore) Get(email string) UserData {
	return s.balances[email]
}

func (s *BalanceStore) Set(email string, u UserData) {
	s.balances[email] = u
}

<<<<<<< HEAD
type StreamCountStore struct {
	mu     sync.Mutex
	path   string
	counts map[string]int64
}

func NewStreamCountStore(path string) *StreamCountStore {
	return &StreamCountStore{path: path}
}

func (s *StreamCountStore) Lock()   { s.mu.Lock() }
func (s *StreamCountStore) Unlock() { s.mu.Unlock() }

func (s *StreamCountStore) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.counts = make(map[string]int64)
			return nil
		}
		return err
	}
	if err := json.Unmarshal(data, &s.counts); err != nil {
		s.counts = make(map[string]int64)
		return nil
	}
	return nil
}

func (s *StreamCountStore) Save() error {
	data, err := json.MarshalIndent(s.counts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *StreamCountStore) AddSeconds(filename string, secs int64) {
	s.counts[filename] += secs
}

func (s *StreamCountStore) GetAll() map[string]int64 {
	result := make(map[string]int64, len(s.counts))
	for k, v := range s.counts {
		result[k] = v
	}
	return result
}

type UniqueUsersStore struct {
	mu    sync.Mutex
	path  string
	users map[string][]string
}

func NewUniqueUsersStore(path string) *UniqueUsersStore {
	return &UniqueUsersStore{path: path}
}

func (s *UniqueUsersStore) Lock()   { s.mu.Lock() }
func (s *UniqueUsersStore) Unlock() { s.mu.Unlock() }

func (s *UniqueUsersStore) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.users = make(map[string][]string)
			return nil
		}
		return err
	}
	if err := json.Unmarshal(data, &s.users); err != nil {
		s.users = make(map[string][]string)
		return nil
	}
	return nil
}

func (s *UniqueUsersStore) Save() error {
	data, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *UniqueUsersStore) Track(email string) {
	if email == "" {
		return
	}
	month := time.Now().Format("2006-01")
	for _, u := range s.users[month] {
		if u == email {
			return
		}
	}
	s.users[month] = append(s.users[month], email)
}

func (s *UniqueUsersStore) GetAllMonths() map[string]int {
	result := make(map[string]int)
	for k, v := range s.users {
		result[k] = len(v)
	}
	return result
}

=======
>>>>>>> origin/master
var (
	configMu   sync.Mutex
	rateLimitMu sync.Mutex
	lastTokenUse = make(map[string]time.Time)
)

func checkTokenRate(email string) bool {
	rateLimitMu.Lock()
	defer rateLimitMu.Unlock()
	last, ok := lastTokenUse[email]
	now := time.Now()
	if ok && now.Sub(last) < 10*time.Second {
		return false
	}
	lastTokenUse[email] = now
	return true
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}

var (
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASS")
	smtpFrom = os.Getenv("SMTP_FROM")
)

<<<<<<< HEAD
var streamSecret []byte

func makeStreamURL(filename string) (string, int64) {
	expiry := time.Now().Add(30 * time.Minute).Unix()
	data := fmt.Sprintf("%s|%d", filename, expiry)
	mac := hmac.New(sha256.New, streamSecret)
	mac.Write([]byte(data))
	sig := hex.EncodeToString(mac.Sum(nil))
	token := base64.RawURLEncoding.EncodeToString([]byte(data + "|" + sig))
	return "/stream/" + token + "/" + filename, expiry
}

func validateStreamToken(token string) (filename string, ok bool) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return
	}
	parts := strings.SplitN(string(raw), "|", 3)
	if len(parts) != 3 {
		return
	}
	filename, expiryStr, sig := parts[0], parts[1], parts[2]

	data := fmt.Sprintf("%s|%s", filename, expiryStr)
	mac := hmac.New(sha256.New, streamSecret)
	mac.Write([]byte(data))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return
	}

	expiryUnix, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil || time.Now().Unix() > expiryUnix {
		return
	}

	ok = true
	return
}

func getEmailFromCookie(r *http.Request) string {
	c, err := r.Cookie("email")
	if err != nil {
		return ""
	}
	return c.Value
}

=======
>>>>>>> origin/master
type MagicToken struct {
	Email   string `json:"email"`
	Expires int64  `json:"expires"`
}

var (
	magicMu    sync.Mutex
	magicPath  = "magic_tokens.json"
)

func loadMagicTokens() (map[string]MagicToken, error) {
	data, err := os.ReadFile(magicPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]MagicToken), nil
		}
		return nil, err
	}
	tokens := make(map[string]MagicToken)
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

func saveMagicTokens(tokens map[string]MagicToken) error {
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(magicPath, data, 0644)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

type loginAuth struct {
	user, pass string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.user), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return []byte(a.pass), nil
	}
	return nil, nil
}

func sendMagicLinkEmail(email, link string) error {
	msg := []byte("From: " + smtpFrom + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: Sign in to nickmonaco.world\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"Click the link below to sign in:\r\n" +
		"\r\n" +
		link + "\r\n" +
		"\r\n" +
		"This link expires in 15 minutes. If you did not request this, you can ignore this email.\r\n")
	log.Printf("connecting to %s:%s as %s", smtpHost, smtpPort, smtpUser)
	addr := smtpHost + ":" + smtpPort

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err := smtp.SendMail(addr, auth, smtpFrom, []string{email}, msg)
	if err == nil {
		return nil
	}
	log.Printf("PLAIN auth failed, trying LOGIN: %v", err)

	auth = &loginAuth{smtpUser, smtpPass}
	return smtp.SendMail(addr, auth, smtpFrom, []string{email}, msg)
}

func main() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	editSecret := os.Getenv("EDIT_SECRET")
	siteDomain := os.Getenv("SITE_DOMAIN")
	if siteDomain == "" {
		siteDomain = "http://localhost:8080"
	}

<<<<<<< HEAD
	key := os.Getenv("STREAM_SECRET")
	if key == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err == nil {
			streamSecret = b
		} else {
			streamSecret = []byte("default-insecure-key-change-me")
		}
	} else {
		streamSecret = []byte(key)
	}

	ssoAggregator := os.Getenv("SSO_AGGREGATOR_URL")
	if ssoAggregator == "" {
		ssoAggregator = "http://localhost:9090"
	}

	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("failed to get executable path: %v", err)
	}
	baseDir := filepath.Dir(exe)
	if err := os.Chdir(baseDir); err != nil {
		log.Fatalf("failed to chdir to %s: %v", baseDir, err)
	}

	broker := NewSSEBroker()
	store := NewBalanceStore("balances.json")
	streamStore := NewStreamCountStore("stream_counts.json")
	uniqueStore := NewUniqueUsersStore("unique_users.json")

	http.HandleFunc("/", indexHandler(editSecret, streamStore, store, siteDomain, ssoAggregator))
	http.HandleFunc("/events", eventsHandler(broker))
	http.HandleFunc("/save-config", saveConfigHandler(editSecret, broker))
	http.Handle("/audio/", audioHandler(siteDomain, ssoAggregator))
	http.Handle("/stream/", streamHandler(siteDomain, ssoAggregator))
	http.HandleFunc("/api/stream-url", streamURLHandler())
	http.Handle("/images/", imagesHandler(siteDomain, ssoAggregator))
	http.HandleFunc("/api/public/whoami", whoamiHandler)
	http.HandleFunc("/api/public/token-supply", tokenSupplyHandler)
	http.HandleFunc("/api/public/total-seconds", totalSecondsHandler)
	http.HandleFunc("/api/public/stream-counts", publicStreamCountsHandler)
	http.HandleFunc("/api/public/report-stream", publicReportStreamHandler(streamStore))
	http.HandleFunc("/api/public/active-users", activeUsersHandler)
	http.HandleFunc("/api/balance", balanceHandler(store))
	http.HandleFunc("/api/use-tokens", useTokensHandler(store, streamStore, uniqueStore))
	http.HandleFunc("/api/stream-counts", streamCountsHandler(streamStore))
	http.HandleFunc("/api/admin/unique-users", uniqueUsersHandler(uniqueStore))
	http.HandleFunc("/sso-login", ssoLoginHandler(ssoAggregator))
=======
	broker := NewSSEBroker()
	store := NewBalanceStore("balances.json")

	http.HandleFunc("/", indexHandler(editSecret))
	http.HandleFunc("/events", eventsHandler(broker))
	http.HandleFunc("/save-config", saveConfigHandler(editSecret, broker))
	http.Handle("/audio/", audioHandler(siteDomain))
	http.Handle("/images/", imagesHandler(siteDomain))
	http.HandleFunc("/api/balance", balanceHandler(store))
	http.HandleFunc("/api/use-tokens", useTokensHandler(store))
>>>>>>> origin/master
	http.HandleFunc("/api/send-magic-link", sendMagicLinkHandler(siteDomain))
	http.HandleFunc("/verify", verifyHandler)
	http.HandleFunc("/buy-tokens", buyTokensHandler(store, siteDomain))
	http.HandleFunc("/success", successHandler(store))
	if webhookSecret != "" {
		http.HandleFunc("/stripe-webhook", stripeWebhookHandler(store, webhookSecret))
	}

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendMagicLinkHandler(domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" || !strings.Contains(body.Email, "@") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "valid email required"})
			return
		}
		token, err := generateToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
			return
		}

		magicMu.Lock()
		tokens, err := loadMagicTokens()
		if err != nil {
			magicMu.Unlock()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
			return
		}
		tokens[token] = MagicToken{Email: body.Email, Expires: time.Now().Add(15 * time.Minute).Unix()}
		saveMagicTokens(tokens)
		magicMu.Unlock()

		if smtpHost == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "SMTP not configured"})
			return
		}
		link := domain + "/verify?token=" + token
		log.Printf("sending magic link to %s via %s:%s", body.Email, smtpHost, smtpPort)
		if err := sendMagicLinkEmail(body.Email, link); err != nil {
			log.Printf("send email error to %s: %v", body.Email, err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to send email"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"ok": "check your email"})
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	magicMu.Lock()
	tokens, err := loadMagicTokens()
	magicMu.Unlock()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mt, ok := tokens[token]
	if !ok || time.Now().Unix() > mt.Expires {
<<<<<<< HEAD
		log.Printf("verify: invalid or expired token %s (ok=%v)", token[:min(8, len(token))], ok)
=======
>>>>>>> origin/master
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Delete used token
	magicMu.Lock()
	delete(tokens, token)
	saveMagicTokens(tokens)
	magicMu.Unlock()

<<<<<<< HEAD
	log.Printf("verify: logged in %s via magic link", mt.Email)

=======
>>>>>>> origin/master
	// Set email cookie and redirect
	http.SetCookie(w, &http.Cookie{
		Name:   "email",
		Value:  mt.Email,
		Path:   "/",
		MaxAge: 365 * 24 * 60 * 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

<<<<<<< HEAD
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func verifySSOTokenWithAggregator(aggregatorURL, token string) (string, bool) {
	if aggregatorURL == "" || token == "" {
		return "", false
	}
	resp, err := http.Get(aggregatorURL + "/api/public/verify-sso?token=" + token)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	var r struct {
		Valid bool   `json:"valid"`
		Email string `json:"email"`
	}
	json.NewDecoder(resp.Body).Decode(&r)
	return r.Email, r.Valid
}

func isAllowedReferer(r *http.Request, domain, aggregatorHost string) bool {
	domainHost := strings.TrimPrefix(domain, "http://")
	domainHost = strings.TrimPrefix(domainHost, "https://")

	// Check Referer header
	if ref := r.Referer(); ref != "" {
		if strings.Contains(ref, domainHost) {
			return true
		}
		if aggregatorHost != "" && strings.Contains(ref, aggregatorHost) {
			return true
		}
	}

	// Check Origin header (more reliable for CORS requests)
	if origin := r.Header.Get("Origin"); origin != "" {
		if strings.Contains(origin, aggregatorHost) {
			return true
		}
	}

	return false
}

func audioHandler(domain, aggregatorURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
=======
func audioHandler(domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ref := r.Referer()
		if ref == "" || !strings.Contains(ref, strings.Split(domain, "://")[1]) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
>>>>>>> origin/master
		filename := strings.TrimPrefix(r.URL.Path, "/audio/")
		if filename == "" || strings.Contains(filename, "..") {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, filepath.Join("audio_files", filename))
	})
}

<<<<<<< HEAD
func streamHandler(domain, aggregatorURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		aggHost := strings.TrimPrefix(aggregatorURL, "http://")
		aggHost = strings.TrimPrefix(aggHost, "https://")
		if !isAllowedReferer(r, domain, aggHost) {
			email := getEmailFromCookie(r)
			if email == "" {
				ssoToken := r.URL.Query().Get("sso_token")
				if ssoToken != "" {
					if e, ok := verifySSOTokenWithAggregator(aggregatorURL, ssoToken); ok {
						email = e
					}
				}
			}
			if email == "" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		path := strings.TrimPrefix(r.URL.Path, "/stream/")
		parts := strings.SplitN(path, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		token, filename := parts[0], parts[1]

		if strings.Contains(filename, "..") {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		fname, ok := validateStreamToken(token)
		if !ok || fname != filename {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		http.ServeFile(w, r, filepath.Join("audio_files", filename))
	})
}

func streamURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		email := getEmailFromCookie(r)
		if email == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "not authenticated"})
			return
		}

		var body struct {
			Filename string `json:"filename"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "filename required"})
			return
		}

		if strings.Contains(body.Filename, "..") {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid filename"})
			return
		}

		streamURL, expiresAt := makeStreamURL(body.Filename)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"url":        streamURL,
			"expires_at": expiresAt,
		})
	}
}

func streamCountsHandler(streamStore *StreamCountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		streamStore.Lock()
		streamStore.Load()
		counts := streamStore.GetAll()
		streamStore.Unlock()
		json.NewEncoder(w).Encode(counts)
	}
}

func uniqueUsersHandler(uniqueStore *UniqueUsersStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		uniqueStore.Lock()
		uniqueStore.Load()
		counts := uniqueStore.GetAllMonths()
		uniqueStore.Unlock()
		json.NewEncoder(w).Encode(counts)
	}
}

func ssoLoginHandler(aggregatorURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("sso_token")
		if token == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}
		// Call aggregator to verify
		resp, err := http.Get(aggregatorURL + "/api/public/verify-sso?token=" + token)
		if err != nil {
			http.Error(w, "verify failed", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		var result struct {
			Valid bool   `json:"valid"`
			Email string `json:"email"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		if !result.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		cookie := &http.Cookie{
			Name:   "email",
			Value:  result.Email,
			Path:   "/",
			MaxAge: 365 * 24 * 60 * 60,
		}
		if r.TLS != nil {
			cookie.SameSite = http.SameSiteNoneMode
			cookie.Secure = true
		}
		http.SetCookie(w, cookie)
		if r.URL.Query().Get("iframe") == "1" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html><body><script>window.parent.postMessage('sso-logged-in','*')</script></body></html>"))
			return
		}
		next := r.URL.Query().Get("redirect")
		if next == "" {
			next = "/"
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}
}

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Vary", "Origin")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	email, err := r.Cookie("email")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"email": ""})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"email": email.Value})
}

func tokenSupplyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	data, err := os.ReadFile("balances.json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]int{"total_supply": 0})
		return
	}
	var balances map[string]struct {
		Balance         int `json:"balance"`
		SecondsListened int `json:"seconds_listened"`
	}
	json.Unmarshal(data, &balances)
	total := 0
	for _, u := range balances {
		total += u.Balance
	}
	json.NewEncoder(w).Encode(map[string]int{"total_supply": total})
}

func activeUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	data, err := os.ReadFile("balances.json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]int{"active_users": 0})
		return
	}
	var balances map[string]struct {
		Balance         int `json:"balance"`
		SecondsListened int `json:"seconds_listened"`
	}
	json.Unmarshal(data, &balances)
	count := 0
	for _, u := range balances {
		if u.SecondsListened > 0 {
			count++
		}
	}
	json.NewEncoder(w).Encode(map[string]int{"active_users": count})
}

func totalSecondsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	data, err := os.ReadFile("balances.json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]int{"total_seconds": 0})
		return
	}
	var balances map[string]struct {
		Balance         int `json:"balance"`
		SecondsListened int `json:"seconds_listened"`
	}
	json.Unmarshal(data, &balances)
	total := 0
	for _, u := range balances {
		total += u.SecondsListened
	}
	// Also include stream counts from public playback
	if scData, err := os.ReadFile("stream_counts.json"); err == nil {
		var sc map[string]int64
		json.Unmarshal(scData, &sc)
		for _, v := range sc {
			total += int(v)
		}
	}
	json.NewEncoder(w).Encode(map[string]int{"total_seconds": total})
}

func publicStreamCountsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCORS(w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	data, err := os.ReadFile("stream_counts.json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]int64{})
		return
	}
	var counts map[string]int64
	json.Unmarshal(data, &counts)
	json.NewEncoder(w).Encode(counts)
}

func setCORS(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")
	}
}

func publicReportStreamHandler(streamStore *StreamCountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		setCORS(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Filename string `json:"filename"`
			Email    string `json:"email"`
			Seconds  int    `json:"seconds"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Filename == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "filename required"})
			return
		}
		if body.Seconds <= 0 {
			body.Seconds = 30
		}
		streamStore.Lock()
		streamStore.Load()
		streamStore.AddSeconds(body.Filename, int64(body.Seconds))
		streamStore.Save()
		streamStore.Unlock()

		// Also update user's seconds_listened if email provided
		if body.Email != "" {
			data, err := os.ReadFile("balances.json")
			if err == nil {
				var balances map[string]struct {
					Balance         int `json:"balance"`
					SecondsListened int `json:"seconds_listened"`
				}
				json.Unmarshal(data, &balances)
				u := balances[body.Email]
				u.SecondsListened += body.Seconds
				balances[body.Email] = u
				out, _ := json.MarshalIndent(balances, "", "  ")
				os.WriteFile("balances.json", out, 0644)
			}
		}

		json.NewEncoder(w).Encode(map[string]string{"ok": "reported"})
	}
}

func imagesHandler(domain, aggregatorURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
=======
func imagesHandler(domain string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ref := r.Referer()
		if ref == "" || !strings.Contains(ref, strings.Split(domain, "://")[1]) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
>>>>>>> origin/master
		filename := strings.TrimPrefix(r.URL.Path, "/images/")
		if filename == "" || strings.Contains(filename, "..") {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, filepath.Join("images", filename))
	})
}

<<<<<<< HEAD
func indexHandler(secret string, streamStore *StreamCountStore, balanceStore *BalanceStore, siteDomain, ssoAggregator string) http.HandlerFunc {
=======
func indexHandler(secret string) http.HandlerFunc {
>>>>>>> origin/master
	return func(w http.ResponseWriter, r *http.Request) {
		configMu.Lock()
		cfg, err := loadConfig()
		configMu.Unlock()
		if err != nil {
			http.Error(w, "Failed to load config", http.StatusInternalServerError)
			log.Printf("config error: %v", err)
			return
		}
		isEdit := secret != "" && r.URL.Query().Get("edit") == secret
		cfgJSON, _ := json.Marshal(cfg)
<<<<<<< HEAD

		signedURLs := make(map[string]string)
		for _, t := range cfg.AudioTracks {
			if !strings.HasPrefix(t.URL, "http") {
				u, _ := makeStreamURL(t.URL)
				signedURLs[t.URL] = u
			}
		}
		for _, a := range cfg.Albums {
			for _, t := range a.Tracks {
				if !strings.HasPrefix(t.URL, "http") {
					u, _ := makeStreamURL(t.URL)
					signedURLs[t.URL] = u
				}
			}
		}
		signedJSON, _ := json.Marshal(signedURLs)

		streamStore.Lock()
		streamStore.Load()
		sc := streamStore.GetAll()
		streamStore.Unlock()
		scJSON, _ := json.Marshal(sc)

		initialSecs := 0
		if email := getEmailFromCookie(r); email != "" {
			balanceStore.Lock()
			balanceStore.Load()
			initialSecs = balanceStore.Get(email).SecondsListened
			balanceStore.Unlock()
		}

		tmpl := template.Must(template.ParseFiles("tmpl/index.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Config":         cfg,
			"ConfigJSON":     template.JS(cfgJSON),
			"IsEdit":         isEdit,
			"EditSecret":     secret,
			"SignedURLs":     signedURLs,
			"SignedURLsJS":   template.JS(signedJSON),
			"StreamCounts":   sc,
			"StreamCountsJS": template.JS(scJSON),
			"InitialSecs":    initialSecs,
=======
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, map[string]interface{}{
			"Config":     cfg,
			"ConfigJSON": template.JS(cfgJSON),
			"IsEdit":     isEdit,
			"EditSecret": secret,
>>>>>>> origin/master
		})
	}
}

func eventsHandler(broker *SSEBroker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := broker.Subscribe()
		defer broker.Unsubscribe(ch)

		configMu.Lock()
		cfg, err := loadConfig()
		configMu.Unlock()
		if err == nil {
			data, _ := json.Marshal(cfg)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}

		for {
			select {
			case msg := <-ch:
				fmt.Fprintf(w, "data: %s\n\n", msg)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}

func saveConfigHandler(secret string, broker *SSEBroker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("edit") != secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		var cfg Config
		if err := json.Unmarshal(bodyBytes, &cfg); err != nil {
			log.Printf("save-config JSON error: %v | body: %s", err, string(bodyBytes))
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		configMu.Lock()
		err := saveConfig(&cfg)
		configMu.Unlock()
		if err != nil {
			http.Error(w, "Failed to save", http.StatusInternalServerError)
			log.Printf("save error: %v", err)
			return
		}
		data, _ := json.Marshal(cfg)
		broker.Broadcast(data)
		w.WriteHeader(http.StatusOK)
	}
}

func balanceHandler(store *BalanceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "email required", http.StatusBadRequest)
			return
		}
		store.Lock()
		if err := store.Load(); err != nil {
			store.Unlock()
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		user := store.Get(email)
		store.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"email":            email,
			"balance":          user.Balance,
			"seconds_listened": user.SecondsListened,
		})
	}
}

<<<<<<< HEAD
func useTokensHandler(store *BalanceStore, streamStore *StreamCountStore, uniqueStore *UniqueUsersStore) http.HandlerFunc {
=======
func useTokensHandler(store *BalanceStore) http.HandlerFunc {
>>>>>>> origin/master
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Email  string `json:"email"`
			Amount int    `json:"amount"`
<<<<<<< HEAD
			Song   string `json:"song"`
=======
>>>>>>> origin/master
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" || body.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}
<<<<<<< HEAD
		uniqueStore.Lock()
		uniqueStore.Load()
		uniqueStore.Track(body.Email)
		uniqueStore.Save()
		uniqueStore.Unlock()

		streamStore.Lock()
		streamStore.Load()
		if body.Song != "" {
			streamStore.AddSeconds(body.Song, 30)
		}
		streamStore.Save()
		allCounts := streamStore.GetAll()
		streamStore.Unlock()

		if !checkTokenRate(body.Email) {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":         "rate limited",
				"stream_counts": allCounts,
			})
			return
		}

		store.Lock()
		if err := store.Load(); err != nil {
			store.Unlock()
=======
		if !checkTokenRate(body.Email) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "rate limited"})
			return
		}
		store.Lock()
		defer store.Unlock()
		if err := store.Load(); err != nil {
>>>>>>> origin/master
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "server error"})
			return
		}
		user := store.Get(body.Email)
		allowedSeconds := user.Balance * 3
		if user.SecondsListened+30 > allowedSeconds {
<<<<<<< HEAD
			store.Unlock()
=======
>>>>>>> origin/master
			json.NewEncoder(w).Encode(map[string]interface{}{
				"balance":          user.Balance,
				"seconds_listened": user.SecondsListened,
				"ok":               false,
<<<<<<< HEAD
				"stream_counts":    allCounts,
=======
>>>>>>> origin/master
			})
			return
		}
		user.SecondsListened += 30
		store.Set(body.Email, user)
		store.Save()
<<<<<<< HEAD
		store.Unlock()

=======
>>>>>>> origin/master
		json.NewEncoder(w).Encode(map[string]interface{}{
			"balance":          user.Balance,
			"seconds_listened": user.SecondsListened,
			"ok":               true,
<<<<<<< HEAD
			"stream_counts":    allCounts,
=======
>>>>>>> origin/master
		})
	}
}

func buyTokensHandler(store *BalanceStore, domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if stripe.Key == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Stripe not configured"})
			return
		}

		var body struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "email required"})
			return
		}

		params := &stripe.CheckoutSessionParams{
			Mode:          stripe.String(string(stripe.CheckoutSessionModePayment)),
			SuccessURL:    stripe.String(domain + "/success?session_id={CHECKOUT_SESSION_ID}"),
			CancelURL:     stripe.String(domain + "/"),
			CustomerEmail: stripe.String(body.Email),
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
						Currency: stripe.String("usd"),
						ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
							Name: stripe.String("500 Tokens"),
						},
						UnitAmount: stripe.Int64(500),
					},
					Quantity: stripe.Int64(1),
				},
			},
		}

		s, err := session.New(params)
		if err != nil {
			log.Printf("stripe session error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Stripe error"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
	}
}

func successHandler(store *BalanceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		s, err := session.Get(sessionID, nil)
		if err != nil {
			log.Printf("stripe session get error: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if s.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if s.CustomerEmail == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		store.Lock()
		if err := store.Load(); err != nil {
			store.Unlock()
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		user := store.Get(s.CustomerEmail)
		user.Balance += 500
		store.Set(s.CustomerEmail, user)
		store.Save()
		balance := user.Balance
		secondsListened := user.SecondsListened
		store.Unlock()

<<<<<<< HEAD
		tmpl := template.Must(template.ParseFiles("tmpl/success.html"))
=======
		tmpl := template.Must(template.ParseFiles("templates/success.html"))
>>>>>>> origin/master
		tmpl.Execute(w, map[string]interface{}{
			"Email":            s.CustomerEmail,
			"Balance":          balance,
			"SecondsListened":  secondsListened,
		})
	}
}

func stripeWebhookHandler(store *BalanceStore, webhookSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusServiceUnavailable)
			return
		}

		sigHeader := r.Header.Get("Stripe-Signature")
		event, err := webhook.ConstructEvent(payload, sigHeader, webhookSecret)
		if err != nil {
			log.Printf("webhook signature error: %v", err)
			http.Error(w, "Signature verification failed", http.StatusBadRequest)
			return
		}

		if event.Type == "checkout.session.completed" {
			var s stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &s); err != nil {
				log.Printf("webhook parse error: %v", err)
				http.Error(w, "Error parsing session", http.StatusBadRequest)
				return
			}

			if s.CustomerEmail != "" {
				store.Lock()
				if err := store.Load(); err != nil {
					store.Unlock()
					http.Error(w, "Error", http.StatusInternalServerError)
					return
				}
				user := store.Get(s.CustomerEmail)
				user.Balance += 500
				store.Set(s.CustomerEmail, user)
				store.Save()
				store.Unlock()
				log.Printf("webhook: credited 500 tokens to %s", s.CustomerEmail)
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

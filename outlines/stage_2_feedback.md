Viewed stage_2.md:1-150
Ran command: `find . -maxdepth 3 -not -path '*/.*'`
Ran command: `find . -maxdepth 3 -not -path '*/.*'`
Viewed main.go:1-30
Viewed client.go:1-56
Viewed credentials.go:1-53
Viewed token.go:1-104
Viewed login.go:1-2
Viewed token_server.go:1-57
Viewed models.go:1-13
Viewed redirect.html:1-16
Ran command: `go build ./...`
Viewed project_overview.md:1-40
Viewed credentials.json:1-12

Your implementation for **Stage 2** is in great shape! You have already completed the core OAuth2 flow: reading credentials, setting up the local HTTP callback server to capture the authorization code, caching tokens with secure `0600` permissions, and initializing the Google Drive client.

Here is a detailed review of how your code is looking, including what is working well, a few critical concurrency/server edge cases to fix, and polish suggestions.

---

### What's Looking Great

1. **Complete Flow Implementation**:
   - [`LoadCredentials`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/credentials.go#L13-L44) correctly parses OAuth credentials via `google.ConfigFromJSON`.
   - [`setClient`](file:///home/JayBlue/personal_projects/google_driver/client.go#L13-L39) implements the token lifecycle properly: checks local cache first via [`ReadTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L17-L41), falls back to interactive OAuth authorization when needed, and initializes the HTTP client.
   - [`setService`](file:///home/JayBlue/personal_projects/google_driver/client.go#L41-L55) successfully initializes `*drive.Service` with the authenticated client.
2. **Security & Permissions**:
   - [`WriteTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L43-L70) sets `0700` for the configuration directory and `0600` for `token.json`, adhering to least-privilege security best practices for secret storage.
3. **Local Callback Server (Approach B)**:
   - Implementing a local server in [`token_server.go`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token_server.go) provides a much better developer and user experience than manual terminal copy-paste.

---

### Critical Issues & Bugs to Fix

#### 1. Deadlock & Data Race in [`GetAuthCode`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token_server.go#L15-L46)
There are two concurrency issues in [`GetAuthCode`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token_server.go#L15-L46):

* **Deadlock if port 8080 is unavailable**: If port 8080 is already in use, `tokenServer.ListenAndServe()` returns an error immediately in the background goroutine. The main goroutine continues to wait on `<-ah.exitChan`, causing the CLI to hang forever.
* **Data race on named return `err`**: `err` is declared in the signature `func GetAuthCode() (code string, err error)` and is written concurrently by the background goroutine (`err = tokenServer.ListenAndServe()`) and the main goroutine (`err = tokenServer.Shutdown(...)`).

**Recommended Fix**:
Use a channel to report server errors and select between the exit trigger and startup errors:

```go
func GetAuthCode() (string, error) {
	ah := authHandler{
		exitChan: make(chan struct{}),
	}
	port := "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", ah.handleGetCode)

	tokenServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Listening for auth code on port %s...", port)
		if err := tokenServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return "", fmt.Errorf("failed to start local server: %w", err)
	case <-ah.exitChan:
		log.Println("Received shutdown trigger")
	}

	if err := tokenServer.Shutdown(context.Background()); err != nil {
		return "", fmt.Errorf("unable to shutdown server: %w", err)
	}

	if ah.code == "" {
		return "", fmt.Errorf("no authorization code received")
	}

	return ah.code, nil
}
```

---

#### 2. Extra Browser Requests (e.g. `/favicon.ico`) in [`handleGetCode`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token_server.go#L48-L56)
When the browser loads `http://localhost:8080/?code=...`, it often immediately fires a second request for `/favicon.ico`. 

In your current code:
```go
func (a *authHandler) handleGetCode(w http.ResponseWriter, r *http.Request) {
    a.code = r.URL.Query().Get("code")
    ...
```
If `/favicon.ico` or any non-root request hits the server:
- `r.URL.Query().Get("code")` will be empty (`""`), which may overwrite the captured code.
- Also, if the user denies access, Google sends `?error=access_denied`, which isn't checked.

**Recommended Fix**:
Filter by path and query parameters before shutting down:
```go
func (a *authHandler) handleGetCode(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(w, "Authorization rejected: "+errMsg, http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	a.code = code
	http.ServeFile(w, r, "./internal/auth/resources/redirect.html")

	go func() {
		a.exitChan <- struct{}{}
	}()
}
```

---

#### 3. Relative Path for `redirect.html` and `credentials.json`
- In [`handleGetCode`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token_server.go#L52), `http.ServeFile(w, r, "./internal/auth/resources/redirect.html")` relies on the current working directory being the root of the project. If the CLI is run from any other folder, this fails with a 404.
  - *Tip*: Use Go 1.16+'s `//go:embed` to embed `redirect.html` into the binary so it works regardless of where the binary is executed.
- Similarly, [`LoadCredentials`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/credentials.go#L14) looks for `./credentials.json`. Consider also checking `~/.config/g-driver-portal/credentials.json` or allowing a path flag in later stages.

---

#### 4. Error Control Flow in [`main`](file:///home/JayBlue/personal_projects/google_driver/main.go#L17-L29)
In [`main.go`](file:///home/JayBlue/personal_projects/google_driver/main.go#L19-L27):
```go
	err := config.setClient()
	if err != nil {
		fmt.Printf("Error setting client: %v\n", err)
	}

	err = config.setService()
	if err != nil {
		fmt.Printf("Error setting service: %v\n", err)
	}
```
If [`setClient`](file:///home/JayBlue/personal_projects/google_driver/client.go#L13) fails, `config.client` remains `nil`. The program continues execution directly into `config.setService()`, which tries to initialize the Drive service with a `nil` HTTP client. 

Consider terminating early on error (e.g. `log.Fatalf` or `return`):
```go
func main() {
	var config apiConfig
	if err := config.setClient(); err != nil {
		log.Fatalf("Error setting client: %v", err)
	}

	if err := config.setService(); err != nil {
		log.Fatalf("Error setting service: %v", err)
	}

	log.Println("Google Drive service successfully initialized!")
}
```

---

### Polish & Clean-Up

1. **Path Construction in [`WriteTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L50-L63)**:
   - Line 14 defines `authTokenPath string = ".config/g-driver-portal/token.json"`.
   - In [`WriteTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L50-L63), `dirPath` is manually created with string concatenation: `homeDir + "/.config/g-driver-portal"`.
   - You can reuse `authTokenPath` with `filepath.Dir(filepath.Join(homeDir, authTokenPath))` and `filepath.Join(homeDir, authTokenPath)`.
2. **Missing Newlines in Print Statements**:
   - [`WriteTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L51) has `fmt.Printf("Directory path: %s", dirPath)` and line 63 has `fmt.Printf("File path: %s", filePath)` without trailing `\n`, which causes them to print together on one line.
3. **Unused Files**:
   - [`internal/auth/login.go`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/login.go) is empty.
   - [`Credentials`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/models.go#L3-L12) in `models.go` is unused since `google.ConfigFromJSON` parses the credentials directly.

---

### Verification Suggestion

To test that Stage 2 works end-to-end, you can add a simple Drive API call in [`main.go`](file:///home/JayBlue/personal_projects/google_driver/main.go) after initializing `cfg.service`:
```go
about, err := config.service.About.Get().Fields("user").Do()
if err != nil {
    log.Fatalf("Drive API check failed: %v", err)
}
fmt.Printf("Authenticated successfully as: %s (%s)\n", about.User.DisplayName, about.User.EmailAddress)
```
This confirms that the token is not only saved and read, but also that Google's API accepts requests from your client.
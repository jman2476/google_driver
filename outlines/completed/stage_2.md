# Stage 2: Authentication & Client Helper Setup

This outline covers the high-level steps required to implement the OAuth 2.0 client helper and token management. Fill in the specific implementation details, code patterns, and design decisions under each section.

---

### 1. Load Application Credentials
* Read and parse `credentials.json` into an `oauth2.Config` object with the appropriate Drive scopes.
* **Details & Implementation Notes**:
  * Use oauth2/google package to generate config
  * 

---

### 2. Check for Cached User Token
* Define a secure local file path (e.g., `token.json` or within `~/.config/`) to store authorization tokens.
* Attempt to read and deserialize existing token data before prompting the user.
* **Details & Implementation Notes**:
  * 
  * 

---

### 3. Obtain Authorization Code (First-Time Login)
* If no cached token exists (or if invalid), generate Google's OAuth consent URL requesting offline access (`access_type=offline`).
* Direct the user to the consent URL and capture the authorization code (e.g., via terminal prompt or temporary localhost callback).
* **Details & Implementation Notes**:
  * Currently will use a copy-paste to capture the authorization code:
Approach A: Terminal Input (Simplest)
Print the URL to the terminal, and prompt the user to paste the code back in.

Methods & Functions to use:

fmt.Printf(...) / fmt.Println(...) to display instructions and the URL.
bufio.NewReader(os.Stdin).ReadString('\n') or fmt.Scan(&authCode) to capture the code from user input.
strings.TrimSpace(authCode) to strip trailing newlines and whitespace.

  * In future, will spin up server to capture the auth code: 
  Approach B: Temporary Local Web Server (Best UX)
Desktop client IDs on Google Cloud allow redirects to http://localhost:<port>. You spin up a temporary local HTTP server, the browser redirects to it automatically upon user consent, and the server grabs the code from the query parameters.

Methods & Functions to use:

net.Listen("tcp", "localhost:8080") or http.Server
r.URL.Query().Get("code") inside your http.HandlerFunc to extract ?code=...
server.Shutdown(ctx) once the code is captured.

Full Gemini output:
Viewed stage_2.md:15-30
Viewed token.go:1-38
Viewed login.go:1-2
Viewed credentials.go:1-53

To implement **Stage 2, Step 3** ("Obtain Authorization Code"), you will primarily use methods from the **`golang.org/x/oauth2`** package and the standard library (`fmt`, `bufio`/`os`, or `net/http`).

Here is a breakdown of the specific methods and two common patterns for capturing the code.

---

### 1. Generating the Consent URL

Use the `AuthCodeURL` method on your `*oauth2.Config`:

```go
authURL := config.AuthCodeURL(
    "state-token",               // state: anti-CSRF token (or any identifier string)
    oauth2.AccessTypeOffline,    // CRUCIAL: requests a refresh token so user stays logged in
    oauth2.ApprovalForce,        // (Optional) forces consent prompt to ensure refresh token is issued
)
```

* **`config.AuthCodeURL(...)`**: Generates the Google OAuth 2.0 URL containing your client ID, requested scopes, and redirect URI.
* **`oauth2.AccessTypeOffline`**: Instructs Google to issue a **refresh token** in addition to a short-lived access token.

---

### 2. Capturing the Authorization Code

There are two approaches to capture the code from the user:

#### Approach A: Terminal Input (Simplest)
Print the URL to the terminal, and prompt the user to paste the code back in.

**Methods & Functions to use:**
* `fmt.Printf(...)` / `fmt.Println(...)` to display instructions and the URL.
* `bufio.NewReader(os.Stdin).ReadString('\n')` or `fmt.Scan(&authCode)` to capture the code from user input.
* `strings.TrimSpace(authCode)` to strip trailing newlines and whitespace.

```go
fmt.Printf("Go to this URL in your browser:\n\n%v\n\n", authURL)
fmt.Print("Enter authorization code: ")

var code string
if _, err := fmt.Scan(&code); err != nil {
    return nil, fmt.Errorf("unable to read authorization code: %w", err)
}
```

---

#### Approach B: Temporary Local Web Server (Best UX)
Desktop client IDs on Google Cloud allow redirects to `http://localhost:<port>`. You spin up a temporary local HTTP server, the browser redirects to it automatically upon user consent, and the server grabs the code from the query parameters.

**Methods & Functions to use:**
* `net.Listen("tcp", "localhost:8080")` or `http.Server`
* `r.URL.Query().Get("code")` inside your `http.HandlerFunc` to extract `?code=...`
* `server.Shutdown(ctx)` once the code is captured.

---

### 3. Exchanging Code for Token (Step 4 preview)

Once you have the authorization code string, exchange it for the actual `*oauth2.Token`:

```go
token, err := config.Exchange(context.Background(), code)
if err != nil {
    return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
}
```

* **`config.Exchange(ctx, code)`**: Sends a POST request to Google's token endpoint (`https://oauth2.googleapis.com/token`) exchanging the one-time authorization code for an `*oauth2.Token` containing:
  * `token.AccessToken`
  * `token.RefreshToken`
  * `token.Expiry`

---

> [!TIP]
> In [`token.go`](file:///home/jeremy/projects/google_driver/internal/auth/token.go#L12), you have `authTokenPath = "~/.config/..."`. Keep in mind that standard Go functions like `os.ReadFile` **do not expand `~`**. Use `os.UserHomeDir()` and `filepath.Join(home, ".config", ...)` to resolve the home directory reliably across platforms.


---

### 4. Exchange Code for Token & Persist
* Exchange the authorization code for an `oauth2.Token` (contains access token and refresh token).
* Save the token to the local cache file with restrictive file permissions (`0600`).
* **Details & Implementation Notes**:
  * 
  * 

---

### 5. Construct Authorized Drive Service Client
* Create an `*http.Client` managed by the OAuth2 config to automatically handle token refreshing.
* Initialize and return the `*drive.Service` instance using `google.golang.org/api/drive/v3`.
* **Details & Implementation Notes**:
  * 
  * 

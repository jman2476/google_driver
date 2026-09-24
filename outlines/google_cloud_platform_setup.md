# Google Cloud Platform (GCP) Setup Guide

This guide walks through configuring Google Cloud Platform (GCP) to enable OAuth 2.0 access for the Google Drive CLI tool.

---

## 1. Create or Select a GCP Project **[Done]**

1. Navigate to the [Google Cloud Console](https://console.cloud.google.com/).
2. Log in with the Google account you want to use.
3. Click the **Project Selector** dropdown at the top-left of the console (next to "Google Cloud").
4. Click **New Project** in the upper-right corner of the modal.
5. Enter a project name (e.g., `google-driver-cli`).
6. *(Optional)* Select an organization or folder if you are using Google Workspace.
7. Click **Create** and wait a few seconds. Ensure your newly created project is selected in the top bar.

---

## 2. Enable the Google Drive API  **[Done]**

1. In the left navigation menu (or search bar), go to **APIs & Services > Library**.
2. Search for **Google Drive API**.
3. Select **Google Drive API** from the results.
4. Click **Enable**.

---

## 3. Configure the OAuth Consent Screen  **[Done]**

Google requires you to set up an OAuth consent screen so users know what application is requesting access and what permissions are required.

1. In the left sidebar, navigate to **APIs & Services > OAuth consent screen**.
2. Choose the **User Type**:
   * **Internal**: Available only if you have a Google Workspace organization and only users inside your organization will use the tool.
   * **External**: Choose this if you are using a standard `@gmail.com` account or want anyone with a Google account to be able to authorize (app will start in "Testing" mode).
3. Click **Create**.

### App Information   **[Done]**
* **App name**: `Google Driver CLI` (or your preferred name).
* **User support email**: Select your email address from the dropdown.
* **Developer contact information**: Enter your email address.
* Leave logos and domain fields blank for now.
* Click **Save and Continue**.

### Scopes   **[Done]**
1. Click **Add or Remove Scopes**.
2. Filter or search for `Google Drive API`.
3. Choose the appropriate scope based on your security preference:
   * **Recommended (Least Privilege)**:
     * `https://www.googleapis.com/auth/drive.file`
     * *Description*: "See, edit, create, and delete only the specific Google Drive files you use with this app."
     * *Why*: The CLI only gets access to files/folders it creates, protecting existing personal drive contents.
   * **Full Access (Broad)**:
     * `https://www.googleapis.com/auth/drive`
     * *Description*: "See, edit, create, and delete all of your Google Drive files."
     * *Why*: Useful if your CLI needs to view, update, or overwrite existing files created outside of the CLI.
4. Check the desired scope, click **Update**, and then click **Save and Continue**.

### Test Users (Crucial for "External" Apps)   **[Done]**
Because your application is in **Testing** status and not officially verified by Google:
1. Click **+ ADD USERS**.
2. Enter the Google/Gmail address(es) that will test the CLI tool.
3. Click **Add**.
4. Click **Save and Continue**, then review the summary and return to the dashboard.

> [!NOTE]
> When your app is in "Testing" mode with an External user type, only listed test users can complete the OAuth login flow. Tokens granted under testing status may expire after 7 days, requiring re-authorization.

---

## 4. Create OAuth 2.0 Client Credentials   **[Done]**

1. In the left sidebar, navigate to **APIs & Services > Credentials**.
2. Click **+ CREATE CREDENTIALS** at the top, then choose **OAuth client ID**.
3. Under **Application type**, select **Desktop app**.
   * *Why Desktop app?* Desktop app credentials allow loopback redirects (`http://localhost:<port>`) or authorization code copy-paste flows without requiring a public HTTPS web endpoint.
4. In the **Name** field, enter `Google Driver CLI Client`.
5. Click **Create**.

---

## 5. Download and Store `credentials.json`   **[Done]**

1. Once created, a popup will display your Client ID and Client Secret. Click **DOWNLOAD JSON** (or find your client under **OAuth 2.0 Client IDs** in the list and click the download button ⬇️ on the right).
2. Rename the downloaded file to `credentials.json`.
3. Place this file where your Go application can locate it (e.g. in the project root during development, or later in `~/.config/google-driver/credentials.json`).

> [!CAUTION]
> **Never commit `credentials.json` or token cache files (`token.json`) to Git.**
> Ensure your `.gitignore` includes:
> ```gitignore
> credentials.json
> token.json
> *.token
> ```

---

## 6. What to Expect on First Login (Consent Screen Warning)   **[]**

When running your CLI for the first time and visiting the generated OAuth consent URL:
1. Google may display a screen titled: **"Google hasn't verified this app"**.
2. Click **Advanced**.
3. Click **Go to Google Driver CLI (unsafe)** to proceed.
4. Review the requested permissions and click **Continue** / **Allow**.

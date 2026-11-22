# Downloading Files from Google Drive

---

### Drive SDK: Which API Call to Use

* **`FilesDownloadCall`** (`service.Files.Download(fileId)`) is a **long-running operation** (LRO) intended for asynchronous Workspace workflows — it does not return a file content stream and is not useful for direct downloads.
* **`FilesGetCall.Download()`** (`service.Files.Get(fileId).Download()`) is the correct method for retrieving raw file bytes — it requests the media endpoint instead of JSON metadata and returns an `*http.Response` with a readable `Body`.
* For native **Google Workspace files** (Docs, Sheets, Slides), the Files.Get media approach will fail with a 403 — these must be **exported** to a standard MIME type (e.g. PDF) via `service.Files.Export()`.

---

### High-Level Download Flow

* Resolve the Drive path into a parent directory ID and filename using existing path utilities (`ValidateDrivePath`).
* Query Drive to locate the specific file and retrieve its file ID and metadata.
* Validate the OAuth token before initiating any network request.
* Use the file ID to request a media stream from the Drive API.
* Determine the local destination path, accounting for whether the target is a directory or an explicit filename.
* Stream the HTTP response body directly to a local file on disk using `io.Copy`.

---

### Key Behaviors & Edge Cases

* The HTTP response body **must always be closed** after use to prevent connection leaks.
* If the user provides a local directory as the destination, the downloaded file should retain its original Drive filename.
* If the user provides an explicit filename as the destination, that name takes precedence.
* Folder downloads require a recursive flag (`-r`) and are a separate concern from single-file downloads.
* Token validation should happen before the download call, consistent with how uploads are handled.
* Download size and progress can be tracked the same way as uploads — `resp.Body` is an `io.Reader` and can be wrapped in a progress tracker.

---

### Workspace Document Exports

* Google Docs, Sheets, and Slides cannot be downloaded as-is — they have no binary blob.
* These file types must be **exported** to a target MIME type (e.g. `application/pdf`, `text/csv`).
* The export API path is distinct from the standard file download path and should be branched on the file's `mimeType`.


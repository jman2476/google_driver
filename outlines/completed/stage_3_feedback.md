# Stage 3 Feedback: Core File Upload & Metadata Management

Your implementation for **Stage 3** is in fantastic shape! You have built a robust, working pipeline that inspects local files, parses folder paths, discovers and creates remote folders on Google Drive, and executes live file uploads. 

All unit and integration tests are currently **passing**, including live Drive API calls (`TestFindFolder`, `TestCreateFolder`, and `TestUploadFile`).

Here is a comprehensive assessment of what is working well, architectural considerations for large files, and polish recommendations as you prepare for Stage 4.

---

### What's Looking Great

1. **Working End-to-End Upload Flow**:
   - [`UploadFile`](file:///home/JayBlue/personal_projects/google_driver/upload.go#L22-L73) connects file validation, destination folder lookup, token validation, and Drive upload execution seamlessly.
   - In [`upload_test.go`](file:///home/JayBlue/personal_projects/google_driver/upload_test.go#L77-L108), `TestUploadFile` successfully uploaded real assets (`bubbletea_gui.gif` and `goose_sqlc-instructions.txt`) to Drive folders with HTTP 200 responses.

2. **Hierarchical Folder Resolution & Creation**:
   - [`ParsePath`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L32-L55) breaks down path strings into structured `DriveFolder` slices.
   - [`FindFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L57-L111) correctly searches folder-by-folder from `"root"` using `'<parentID>' in parents` and `trashed = false`, resolving the earlier panic.
   - [`CreateFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L113-L170) dynamically creates target folders under their respective parents on Drive when needed.

3. **Proactive Token Refresh**:
   - In [`client.go`](file:///home/JayBlue/personal_projects/google_driver/client.go#L57-L78), [`validateToken`](file:///home/JayBlue/personal_projects/google_driver/client.go#L57) calls [`EnsureValidToken`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/refresh.go#L37-L50) before both folder operations and file uploads.
   - If refreshed, it immediately updates `token.json` via [`WriteTokenData`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/token.go#L44-L70), cleanly completing Steps 3 & 4 of the token refresh outline.

4. **Testing Rigor**:
   - You have added thorough tests with table-driven test cases across [`detect_test.go`](file:///home/JayBlue/personal_projects/google_driver/detect_test.go), [`directory_test.go`](file:///home/JayBlue/personal_projects/google_driver/directory_test.go), and [`upload_test.go`](file:///home/JayBlue/personal_projects/google_driver/upload_test.go).

---

### Architectural Considerations & Areas for Improvement

#### 1. Streaming vs. In-Memory Buffering for Large Files
In [`detect.go:48-60`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L48-L60):
```go
data, err := os.ReadFile(file.AbsPath)
// ...
file.Data = data
file.Reader = bytes.NewReader(file.Data)
```
- **The Issue**: `os.ReadFile` reads the **entire file into memory (RAM)**. For small files (a few kilobytes or megabytes), this works fine. However, if a user uploads a 500 MB video (like `gui_example.mov`) or a 4 GB ISO, the process will consume gigabytes of RAM or crash with an out-of-memory error.
- **Recommended Fix**: Stream directly from the disk using `*os.File` instead of buffering the whole payload into memory:
  ```go
  type FileData struct {
      AbsPath string
      Info    os.FileInfo
      Mime    string
      Reader  io.ReadCloser // or *os.File
  }
  ```
  In [`ValidateFile`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L39-L63):
  ```go
  fileHandle, err := os.Open(file.AbsPath)
  if err != nil {
      return FileData{}, fmt.Errorf("unable to open file: %w", err)
  }
  file.Reader = fileHandle
  ```
  The Google Drive Go SDK (`uploadBuilder.Media(file.Reader)`) accepts any `io.Reader` and streams the file in chunks without buffering it all into RAM.

---

#### 2. Robust MIME Type Detection Fallback
In [`detect.go:55-59`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L55-L59):
```go
nameParts := strings.Split(file.Info.Name(), ".")
extension := "." + nameParts[len(nameParts)-1]
mimeType := mime.TypeByExtension(extension)
```
- **Edge Cases**:
  1. If a file has no extension (e.g. `LICENSE`, `Makefile`, `Dockerfile`), `nameParts` has length 1, and `extension` becomes `".LICENSE"`, which `mime.TypeByExtension` cannot identify. Use the standard library function `filepath.Ext(file.Info.Name())` instead.
  2. If `mime.TypeByExtension` returns empty string `""` (common for unknown extensions or files without extensions), `file.Mime` remains empty.
- **Recommended Fix**:
  Add content sniffing (reading the first 512 bytes) and fallback to `application/octet-stream`:
  ```go
  ext := filepath.Ext(file.Info.Name())
  mimeType := mime.TypeByExtension(ext)
  if mimeType == "" {
      // Sniff first 512 bytes
      buf := make([]byte, 512)
      n, _ := fileHandle.Read(buf)
      fileHandle.Seek(0, io.SeekStart) // Rewind back to start for the upload!
      if n > 0 {
          mimeType = http.DetectContentType(buf[:n])
      } else {
          mimeType = "application/octet-stream"
      }
  }
  ```

---

#### 3. Returning Upload Results from `UploadFile`
In [`upload.go:22`](file:///home/JayBlue/personal_projects/google_driver/upload.go#L22):
```go
func (cfg *apiConfig) UploadFile(path string, targetDir string) error
```
- **The Issue**: `UploadFile` currently returns only `error`. Once a file is uploaded, the caller (such as `main.go` or CLI commands) often needs to know the created file's ID, web link (`webViewLink`), or confirm the filename.
- **Recommended Fix**: Return the created `*drive.File` (or custom result struct):
  ```go
  func (cfg *apiConfig) UploadFile(path string, targetDir string) (*drive.File, error)
  ```
  Update line 61 to request the web link:
  ```go
  response, err := uploadBuilder.Fields("id", "name", "parents", "size", "webViewLink").Do()
  ```
  This makes it trivial to display:
  `Upload complete! File ID: 1a2b3c... View at: https://drive.google.com/...`

---

#### 4. Safety in `TestFindFolder`
In [`directory_test.go:93-98`](file:///home/JayBlue/personal_projects/google_driver/directory_test.go#L93-L98):
```go
folder, err := config.FindFolder(cases)
if err != nil {
    t.Logf("Fail: TFF error: %v\n", err)
}

fmt.Printf("Folder: %s\n", folder.Name)
```
- **The Issue**: If `FindFolder` returns an error (e.g. folder deleted or not found), `t.Logf` logs the message but allows execution to continue to line 98. Accessing `folder.Name` when `folder == nil` will panic with a segmentation fault.
- **Recommended Fix**: Use `t.Fatalf` or verify `folder != nil` before reading properties:
  ```go
  folder, err := config.FindFolder(cases)
  if err != nil {
      t.Fatalf("FindFolder failed: %v", err)
  }
  if folder == nil {
      t.Fatalf("Expected folder, got nil")
  }
  ```

---

### Minor Polish & Code Cleanliness

1. **Trailing Newlines in Error Strings**:
   In Go, error messages should not end with `\n` (e.g., `fmt.Errorf("failed to validate file: %v\n", err)` in [`upload.go:26`](file:///home/JayBlue/personal_projects/google_driver/upload.go#L26)). Removing the trailing `\n` keeps CLI logs consistent.
2. **Close File Readers**:
   When opening file handles for upload, make sure file handles are closed using `defer` to prevent file descriptor leaks during bulk uploads.
3. **Empty Folder Deletion**:
   [`DeleteFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L172) is currently a stub returning `nil, nil`. You can wire this to `cfg.service.Files.Delete(folderID).Do()` when ready.

---

### Progress vs. Project Roadmap

| Stage 3 Requirement | Status | Implementation Details |
| :--- | :--- | :--- |
| **Local File Inspection** | **Done** | [`FindFile`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L20), [`ValidateFile`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L39) |
| **MIME Type Detection** | **Done** | [`mime.TypeByExtension`](file:///home/JayBlue/personal_projects/google_driver/detect.go#L57) (sniffing fallback suggested) |
| **Destination Folder Targeting** | **Done** | [`ParsePath`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L32), [`FindFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L57), [`CreateFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L113) |
| **Drive File Upload Execution** | **Done** | [`WriteUploadMetadata`](file:///home/JayBlue/personal_projects/google_driver/upload.go#L9), [`UploadFile`](file:///home/JayBlue/personal_projects/google_driver/upload.go#L22) |
| **Proactive Token Refresh** | **Done** | [`validateToken`](file:///home/JayBlue/personal_projects/google_driver/client.go#L57), [`EnsureValidToken`](file:///home/JayBlue/personal_projects/google_driver/internal/auth/refresh.go#L37) |

### Verdict

**Stage 3 is a success.** You have working, verified upload and folder management logic. Addressing streaming reads for large files and returning upload metadata will set you up perfectly for **Stage 4** (CLI commands, flags with Cobra/CLI, and progress indicators).

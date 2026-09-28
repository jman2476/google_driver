# Stage 3: Core File Upload & Metadata Management

This outline covers the high-level steps required to inspect local files, detect MIME types, execute file uploads to Google Drive, and handle destination folders. Fill in the specific implementation details, code patterns, and design decisions under each section.

---

### 1. Local File Inspection & Validation
* Validate that the target file path exists, is readable, and is a regular file (not a directory).
* Gather basic file attributes (size, filename) required for upload and progress tracking.
* **Details & Implementation Notes**:
  * 
  * 

---

### 2. MIME Type Detection
* Determine the appropriate MIME type so Google Drive correctly identifies and previews the file.
* Implement a detection strategy (e.g., file extension mapping with content sniffing fallback).
* **Details & Implementation Notes**:
  * 
  * 

---

### 3. Destination Folder Targeting
* Allow the upload to target either the root Drive directory or a specific Google Drive folder ID.
* Set the appropriate parent reference on the Drive file metadata before uploading.
* **Details & Implementation Notes**:
  * 
  * 

---

### 4. Execute Drive File Upload
* Construct the `drive.File` metadata object with the filename, MIME type, and parent folder.
* Open the local file reader and stream the file contents to the Google Drive API.
* **Details & Implementation Notes**:
  * 
  * 

---

### 5. Return Upload Metadata & Result
* Capture and return the uploaded file details (e.g., file ID, name, size, web view link).
* Handle upload errors gracefully and verify successful upload.
* **Details & Implementation Notes**:
  * 
  * 


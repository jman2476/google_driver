# Stage 4: CLI Interface, Commands & User Experience Polish

This outline covers the high-level steps required to build a structured command-line interface (CLI), support subcommands and flags, display upload progress indicators, and provide graceful error handling and user feedback. Fill in the specific implementation details, code patterns, and design decisions under each section.

---

### 1. Command Routing & Dispatcher Architecture
* Define a consistent command pattern for all CLI tools (e.g., `mv`/upload, `mkdir`, `rm`, `login`, `logout`, `account`).
* Provide a centralized help command and print clear usage instructions when arguments are missing or invalid.
* Standardize command callbacks and execution lifecycles across the application.
* **Details & Implementation Notes**:
  * 
  * 

---

### 2. Flag & Option Parsing
* Support optional flags for commands (e.g., target destination folder `--dir`/`-d`, overwrite confirmations, or verbose logging).
* Handle optional arguments cleanly without strict reliance on fixed positional indexes.
* Support environment variables or configuration flags for custom credential or cache paths.
* **Details & Implementation Notes**:
  * 
  * 

---

### 3. Upload Progress Tracking
* Wrap the file upload stream (`io.Reader`) to monitor bytes transferred in real time against total file size.
* Display dynamic terminal feedback (e.g., progress bar, percentage counter, transfer rate) during large file uploads.
* Ensure progress output cleanly finishes before final success links are printed.
* **Details & Implementation Notes**:
  * 
  * 

---

### 4. Remote Drive Management Commands
* Implement remaining core Drive operations to round out the toolset (e.g., `rm` for deleting or trashing files, `ls` for listing remote folder contents).
* Include confirmation prompts or dry-run checks before executing destructive remote actions.
* Handle remote path-to-ID resolution for all file and folder management commands.
* **Details & Implementation Notes**:
  * 
  * 

---

### 5. Error Handling, Exit Codes & CLI Polish
* Translate low-level API and network errors into clear, actionable terminal messages (e.g., rate limits, invalid grants, network interruptions).
* Standardize Unix exit codes across all commands (`0` on success, distinct non-zero codes for errors).
* Provide clean, readable success summaries with clickable web view links upon operation completion.
* **Details & Implementation Notes**:
  * 
  * 

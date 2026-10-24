Here is a comprehensive progress check on your **[Stage 4: CLI Interface, Commands & User Experience Polish](file:///home/jeremy/projects/google_driver/outlines/stage_4.md)** checklist.

---

### Progress Summary Table

| Stage 4 Requirement | Status | Current Implementation Details |
| :--- | :---: | :--- |
| **1. Command Routing & Dispatcher Architecture** | **Partially Done** | Registered commands in [`commandRegistry`](file:///home/jeremy/projects/google_driver/command.go#L11); missing `help` command, usage output, and uniform execution lifecycle. |
| **2. Flag & Option Parsing** | **Not Started** | [`parseArgs`](file:///home/jeremy/projects/google_driver/parse.go#L12) strictly checks positional slices; no flags (`--dir`, `-d`), option parser, or env var support. |
| **3. Upload Progress Tracking** | **Done** (Minor Polish) | [`Tracker`](file:///home/jeremy/projects/google_driver/internal/progress/reader.go#L13) wraps stream, prints dynamic percentage, transfer rate, ETA, and bar. Passes all tests. |
| **4. Remote Drive Management Commands** | **Mostly Done** | `ls`, `rm`, `mv`, `mkdir`, `account`, `login`, and `logout` are functional with path-to-ID resolution. Missing confirmation prompts / dry-run for `rm`. |
| **5. Error Handling, Exit Codes & Polish** | **Partially Done** | Uses `0` / `1` exits; web links printed for `upload` and `mkdir`. Missing API error translation, distinct exit codes, and deletion feedback. |

---

### Detailed Review by Checklist Item

#### 1. Command Routing & Dispatcher Architecture
* **What's Working:**
  * Defined [`cliCommand`](file:///home/jeremy/projects/google_driver/command.go#L5) and centralized commands in [`commandRegistry`](file:///home/jeremy/projects/google_driver/command.go#L11) (`ls`, `rm`, `mv`, `mkdir`, `login`, `logout`, `account`).
* **Needs Work:**
  * **Centralized Help & Usage:** There is no `help` command or `--help` flag. When arguments are missing or invalid, [main.go](file:///home/jeremy/projects/google_driver/main.go#L28-L31) prints `Argument error: not enough arguments` and exits without displaying usage instructions or available commands.
  * **Nil Callback Panic on Unknown Tool:** If a user enters an unrecognized command (e.g. `g-driver foo`), [`commandRegistry()[args.tool]`](file:///home/jeremy/projects/google_driver/main.go#L64) returns a zero-value struct with a `nil` callback. Calling `tool.callback(...)` on [line 65](file:///home/jeremy/projects/google_driver/main.go#L65) triggers a runtime panic. Check if the tool exists in the map first.
  * **Execution Lifecycle Uniformity:** In [main.go](file:///home/jeremy/projects/google_driver/main.go#L33-L50), `login` and `logout` are hardcoded into an initial `switch` statement that duplicates execution outside the registry, while all other commands go through the registry callback.

---

#### 2. Flag & Option Parsing
* **Current State:**
  * [`parseArgs`](file:///home/jeremy/projects/google_driver/parse.go#L12) currently uses rigid positional argument slicing (`len(args)`) mapping to `tool`, `source`, and `target`.
* **Needs Work:**
  * **Flag Support:** No support for flags like `--dir`/`-d`, `--dry-run`, `--yes`/`-y`, or `--verbose`.
  * **Flexible Positional Ordering:** Commands like `mv <file> --dir <folder>` or `mv --dir <folder> <file>` are rejected as `too many arguments`.
  * **Config / Env Var Overrides:** `credentials.json` and `token.json` paths remain hardcoded in [credentials.go](file:///home/jeremy/projects/google_driver/internal/auth/credentials.go#L14) and [token.go](file:///home/jeremy/projects/google_driver/internal/auth/token.go#L14) without flag or environment variable override options (`GOOGLE_APPLICATION_CREDENTIALS` or `--credentials`).
  * *Tip:* You can use standard library [`flag`](https://pkg.go.dev/flag) (with `flag.NewFlagSet`) or custom parsing to read optional flags.

---

#### 3. Upload Progress Tracking
* **What's Working:**
  * Real-time stream monitoring implemented via [`Tracker.Read`](file:///home/jeremy/projects/google_driver/internal/progress/reader.go#L30) wrapping `*os.File`.
  * Terminal-width-aware bar formatting ([`PrintProgressBar`](file:///home/jeremy/projects/google_driver/internal/progress/reader.go#L59)), human-readable transfer rate ([`formatRate`](file:///home/jeremy/projects/google_driver/internal/progress/reader.go#L109)), and time remaining calculation ([`calcTimeRemaining`](file:///home/jeremy/projects/google_driver/internal/progress/reader.go#L143)).
  * Fully tested in [`reader_test.go`](file:///home/jeremy/projects/google_driver/internal/progress/reader_test.go).
* **Minor Polish:**
  * `PrintProgressBar` updates the line with `\r` and no trailing newline. When the upload completes, [upload.go:79](file:///home/jeremy/projects/google_driver/upload.go#L79) prints `File uploaded to ...`, which directly overwrites the progress bar on terminal screens. Printing a `\n` upon completion or EOF will keep both the completed bar and the URL visible.

---

#### 4. Remote Drive Management Commands
* **What's Working:**
  * Core commands implemented:
    * `ls`: [`commandList`](file:///home/jeremy/projects/google_driver/command_list.go#L10) and [`listFolder`](file:///home/jeremy/projects/google_driver/directory.go#L171) list contents with formatted MIME types and file sizes.
    * `rm`: [`commandRemove`](file:///home/jeremy/projects/google_driver/command_remove.go#L9) and [`removeFile`](file:///home/jeremy/projects/google_driver/remove.go#L17) resolve parents and delete files.
    * `mkdir`: [`createFolder`](file:///home/jeremy/projects/google_driver/directory.go#L112) resolves hierarchical folders.
    * `mv`: [`uploadFile`](file:///home/jeremy/projects/google_driver/upload.go#L24) streams and uploads files with metadata.
  * Integration tests for removal and error handling ([`remove_test.go`](file:///home/jeremy/projects/google_driver/remove_test.go)) pass cleanly.
* **Needs Work:**
  * **Confirmation / Dry-Run:** [`removeFile`](file:///home/jeremy/projects/google_driver/remove.go#L52) immediately invokes `cfg.service.Files.Delete(id).Do()` (which permanently deletes rather than moves to trash). Adding an interactive confirmation (`Are you sure? [y/N]`) or a `--dry-run` / `--yes` flag will prevent accidental deletions.

---

#### 5. Error Handling, Exit Codes & CLI Polish
* **What's Working:**
  * Clean exit codes (`os.Exit(0)` and `os.Exit(1)`) are used in [main.go](file:///home/jeremy/projects/google_driver/main.go).
  * Web view links are printed for uploaded files and created folders.
* **Needs Work:**
  * **Silent Removal:** In [`remove.go:58`](file:///home/jeremy/projects/google_driver/remove.go#L58), `removeFile` returns `nil, nil`. In `main.go`, `if file != nil` is skipped, so successful file deletions happen completely silently without user confirmation output. Returning a populated `*drive.File{Name: fileName}` allows `main.go` to print `Performed 'Remove' on 'filename'`.
  * **API Error Translation:** Raw API/OAuth errors (e.g. expired tokens, rate limit HTTP 429, 403 Forbidden, invalid credentials) are printed as raw Go error strings without friendly guidance (e.g. suggesting `g-driver login`).
  * **Standardized Exit Codes:** Exit codes are currently only binary `0` and `1`. You could define specific codes (e.g., `2` for syntax/argument errors, `3` for authentication errors, `4` for not-found errors).

---

### Suggested Next Steps
1. Add an unknown-command guard and a `help` command in [command.go](file:///home/jeremy/projects/google_driver/command.go) / [main.go](file:///home/jeremy/projects/google_driver/main.go) to print usage text.
2. Upgrade [parse.go](file:///home/jeremy/projects/google_driver/parse.go) to support flags (especially for optional destination directory `--dir` and confirmation `--yes`).
3. Return the deleted file metadata from [`removeFile`](file:///home/jeremy/projects/google_driver/remove.go#L58) and add a user confirmation prompt before deletion.

Potato
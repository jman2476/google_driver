# Flag & Option Parsing Architecture

This outline covers the high-level steps required to transition from fixed positional argument parsing to flexible flag and option parsing. It details how to handle global flags, command-specific flags, optional arguments, and configuration overrides without breaking existing command dispatching. Fill in the specific implementation details, patterns, and design decisions under each section.

---

### 1. Requirements & Strategy for Flag Parsing
* Identify which commands require flags versus positional arguments (e.g., optional target directories `--dir`/`-d` for `mv`, confirmation `--yes`/`-y` or `--dry-run` for `rm`, verbose mode `-v`/`--verbose`).
* Decide on the parsing approach (leveraging Go's standard library `flag` package with sub-command `FlagSet`s, or building a lightweight custom token scanner to preserve current simple semantics).
* Determine support for short flags (`-d`), long flags (`--dir`), boolean switches (`--yes`), and flags with values (`--dir /path` or `--dir=/path`).
* Establish conventions for flag placement, allowing flags to appear before or after positional arguments.
* **Details & Implementation Notes**:
  * 
  * 

---

### 2. Upgrading the Arguments Model
* Extend the `arguments` struct in `parse.go` to store parsed flags separately from positional arguments (such as destination directory, confirmation booleans, dry-run flags, and config path overrides).
* Maintain backward compatibility for commands that expect positional source and target inputs while flags are being adopted.
* Define sensible default values for optional flags when they are omitted by the user.
* **Details & Implementation Notes**:
  * 
  * 

---

### 3. Parsing Pipeline & Token Separation
* Inspect incoming raw CLI arguments (`os.Args`) to separate the command name (subcommand) from subsequent options and parameters.
* Tokenize and separate flags from non-flag positional arguments so argument order does not cause parsing failures.
* Stop flag parsing cleanly upon encountering positional arguments or standard double-dash delimiters (`--`).
* Return clear, descriptive parsing errors when unexpected flags, malformed flag syntax, or missing flag values are encountered.
* **Details & Implementation Notes**:
  * 
  * 

---

### 4. Command-Specific Flag Binding & Validation
* Bind appropriate flags to their respective commands in `commandRegistry` (for instance, ensuring `--dir` is recognized for `mv`, while `--yes` and `--dry-run` apply to `rm`).
* Validate that mutually exclusive flags or incompatible flag combinations are caught early.
* Ensure commands that take no flags fail gracefully with informative error messages when unexpected flags are passed.
* Propagate parsed flag values into the execution callbacks defined in `cliCommand`.
* **Details & Implementation Notes**:
  * 
  * 

---

### 5. Configuration & Environment Overrides
* Support configuration flags for sensitive or environment-specific paths (such as `--credentials` for credentials JSON location and `--cache` for token storage).
* Implement a clear precedence hierarchy: explicit CLI flags take highest priority, followed by environment variables, falling back to existing default local paths in `internal/auth`.
* Validate file existence and permissions for custom configuration paths before proceeding to Google Drive API initialization.
* **Details & Implementation Notes**:
  * 
  * 


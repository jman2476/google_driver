# Custom CLI Flag & Argument Parser Guide

This document outlines the design, trade-offs, and recommended implementation for replacing fixed positional argument parsing with a custom token scanner in `google_driver`.

---

## 1. Architectural Assessment: `internal/argparser` vs. `package main`

### Why `internal/argparser` creates unnecessary friction
In Go, an `internal/` package **cannot import `package main`**.
- In `google_driver`, commands, configuration (`apiConfig`), and command dispatching (`cliCommand`) all live in `package main` (`main.go`, `command.go`, etc.).
- If command-specific parsing (e.g., `parseMove()`, `parseRemove()`) lives in `internal/argparser`, that package cannot reference any structs or callbacks defined in `main`.
- You already have `parse.go` in `package main`. Keeping argument parsing in `parse.go` allows it to directly interface with your CLI data types and error constants without import cycles or unnecessary boundaries.

### Why per-command parser functions are overcomplicating things
Writing individual functions like `parseList()`, `parseMove()`, `parseRemove()`, and `parseMkDir()` creates repetitive boilerplate. Almost all CLI commands share the same structure:
1. **Subcommand / Tool** (`args[1]`, e.g., `mv`, `rm`, `ls`)
2. **Positional Arguments** (e.g., source file, target directory)
3. **Flags / Options** (e.g., `-d`, `--dir`, `-y`, `--yes`)

Instead of writing 8 separate parsers, you only need **one generic tokenizer pass** that separates arguments into:
- `tool` (the subcommand)
- `positional` (`[]string`)
- `flags` (`map[string]string`)

Individual commands simply query positional arguments and flags by name.

---

## 2. Pitfalls in Naive Token Scanning

When building a custom CLI scanner, avoid these four common pitfalls:

1. **Index Out-of-Range Panics**:
   Hardcoding `parsed.Command = args[1]`, `parsed.SourcePath = args[2]`, `parsed.TargetPath = args[3]` assumes every command receives at least 4 arguments. Commands like `google_driver ls` or `google_driver login` only have 2 tokens (`[binary, "ls"]`), causing immediate runtime panics.

2. **Fixed Positions Break Flag Order**:
   Assuming flags only start at `args[4]` breaks common CLI usage like:
   ```bash
   google_driver rm file.txt -y
   ```
   Here, `-y` is at index 3. If index 3 is blindly assigned to `TargetPath`, the flag is never recognized as a flag and positional arguments get corrupted.

3. **`strings.Contains` vs. `strings.HasPrefix`**:
   Using `strings.Contains(arg, "-")` mistakenly flags normal filenames containing hyphens (e.g. `report-2026.pdf` or `work-backup`) as flags. Always use `strings.HasPrefix(token, "-")`.

4. **Boolean Switches vs. Value Flags**:
   Assuming `args[i+1]` is always the value of `args[i]` breaks boolean switches (`-y`, `--yes`, `--dry-run`). It will mistakenly consume the subsequent file path as the flag's value.

---

## 3. Recommended Implementation (`parse.go`)

Upgrade `parse.go` in `package main` with the following clean, order-independent tokenizer:

```go
package main

import (
	"errors"
	"strings"
)

var (
	ErrMissingArgs = errors.New("command required")
)

type arguments struct {
	tool       string
	positional []string
	flags      map[string]string // e.g. "dir": "folder1", "yes": "true"
}

// HasFlag checks if any of the provided flag aliases were set (e.g., "y", "yes")
func (a arguments) HasFlag(names ...string) bool {
	for _, name := range names {
		if _, ok := a.flags[name]; ok {
			return true
		}
	}
	return false
}

// GetFlag returns the value of the first matching flag alias, or ""
func (a arguments) GetFlag(names ...string) string {
	for _, name := range names {
		if val, ok := a.flags[name]; ok {
			return val
		}
	}
	return ""
}

func parseArgs(rawArgs []string) (arguments, error) {
	if len(rawArgs) < 2 {
		return arguments{}, ErrMissingArgs
	}

	parsed := arguments{
		tool:  rawArgs[1],
		flags: make(map[string]string),
	}

	args := rawArgs[2:]
	for i := 0; i < len(args); i++ {
		token := args[i]

		// Check if token is a flag
		if strings.HasPrefix(token, "-") {
			flagName := strings.TrimLeft(token, "-")

			// Check for --flag=value syntax
			if strings.Contains(flagName, "=") {
				parts := strings.SplitN(flagName, "=", 2)
				parsed.flags[parts[0]] = parts[1]
				continue
			}

			// Check if next token is a value or another flag
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				parsed.flags[flagName] = args[i+1]
				i++
			} else {
				// Standalone boolean flag (e.g. -y, --yes)
				parsed.flags[flagName] = "true"
			}
			continue
		}

		// Non-flag arguments are positional
		parsed.positional = append(parsed.positional, token)
	}

	return parsed, nil
}
```

---

## 4. Why This Architecture Works Best

1. **Order Independent**: 
   Both `google_driver mv file.txt -d /folder` and `google_driver mv -d /folder file.txt` yield identical, expected data:
   - `tool`: `"mv"`
   - `positional`: `["file.txt"]`
   - `flags`: `{"d": "/folder"}`

2. **Zero Runtime Panics**:
   Safely processes input slices of any length.

3. **Decoupled Command Handlers**:
   Command functions only look up what they care about:
   ```go
   targetDir := args.GetFlag("dir", "d")
   skipConfirm := args.HasFlag("yes", "y")
   ```

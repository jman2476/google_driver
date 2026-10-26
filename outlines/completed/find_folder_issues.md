Viewed directory_test.go:70-104
Viewed directory.go:1-116
Ran command: `go test -v -run TestFindFolder`
Viewed token.json

This error is a **nil pointer dereference panic** occurring at line 94 of [`directory_test.go`](file:///home/JayBlue/personal_projects/google_driver/directory_test.go#L94):

```text
panic: runtime error: invalid memory address or nil pointer dereference
github.com/jman2476/google-driver.TestFindFolder
    /home/JayBlue/personal_projects/google_driver/directory_test.go:94 +0xc4
```

Here is a breakdown of why this happens and how to fix it:

---

### 1. The Immediate Cause (in `directory_test.go`)

In [`directory_test.go`](file:///home/JayBlue/personal_projects/google_driver/directory_test.go#L92-L94):

```go
folder, err := config.FindFolder(cases)

fmt.Printf("Folder: %s", folder.Name) // <--- Panics here at line 94
```

Because [`FindFolder`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L57-L107) could not find the folder `"/google-driver_testing"`, it returned `folder = nil`.

Attempting to read `folder.Name` when `folder` is `nil` triggers a segmentation violation (`SIGSEGV`). The address `addr=0x170` in the panic message is the internal memory offset of the `Name` field inside the `drive.File` struct.

In your test, you should always check `err` and verify `folder != nil` first:

```go
folder, err := config.FindFolder(cases)
if err != nil {
    t.Fatalf("FindFolder failed: %v", err)
}
if folder == nil {
    t.Fatalf("Expected folder, got nil")
}
fmt.Printf("Folder: %s\n", folder.Name)
```

---

### 2. Why `FindFolder` Returned `(nil, nil)`

In [`directory.go`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L83-L107):

```go
var tracker = struct {
    ParentID    string
    ParentName  string
    CurrentFile *drive.File
}{}

// ... search loop ...

return tracker.CurrentFile, nil
```

If no matching folder is found, `tracker.CurrentFile` is never populated, so the function returns `nil` for the file and `nil` for the error!

You should check if a folder was found, and return an error if it wasn't (using the [`ErrDirNotFound`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L18) you already defined):

```go
if tracker.CurrentFile == nil {
    return nil, ErrDirNotFound
}
return tracker.CurrentFile, nil
```

---

### 3. Why the Folder Wasn't Found in Google Drive

There are three reasons the search is failing to match:

#### A. Does `"google-driver_testing"` exist in Drive?
If you haven't created a folder named `google-driver_testing` in your Google Drive yet, the query returns 0 folders. Also, if your app is authorized with `drive.DriveFileScope`, it can only see folders **created by this app**, not folders created manually through the browser.

#### B. Drive API v3 does not return `parents` by default
In [`directory.go:72`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L72):
```go
folderList, err := cfg.service.Files.List().Q(searchQuery).Do()
```
In Google Drive API v3, `Files.List()` only returns `id`, `name`, and `mimeType`. It **does not return `parents`** unless you explicitly request it:
```go
folderList, err := cfg.service.Files.List().
    Q(searchQuery).
    Fields("files(id, name, parents)").
    Do()
```

#### C. Comparing Names vs. IDs in the Search Loop
In [`directory.go:92-96`](file:///home/JayBlue/personal_projects/google_driver/directory.go#L92-L96):
```go
if folder.Parent == "" && len(f.Parents) == 0 {
    tracker.CurrentFile = f
    break
} else if folder.Parent == f.Parents[0] {
    ...
```
1. **`len(f.Parents) == 0`**: Folders located in your root "My Drive" are **not** empty—they have a parent ID pointing to your Drive root. `len(f.Parents)` is never `0`.
2. **`folder.Parent == f.Parents[0]`**: `folder.Parent` is a human-readable name string (e.g. `"my_folder"`), while `f.Parents[0]` is a Google Drive resource ID (e.g. `"1AbC...9z"`). They will never be equal.

---

### Recommended Folder Search Pattern

To resolve folder paths like `"/parent/child"` in Drive, search step-by-step from root:

```go
func (cfg *apiConfig) FindFolder(path string) (*drive.File, error) {
	folders, err := ParsePath(path)
	if err != nil {
		return nil, fmt.Errorf("path error: %w", err)
	}

	currentParentID := "root"
	var currentFolder *drive.File

	for _, folder := range folders {
		query := fmt.Sprintf(
			"mimeType = '%s' and name = '%s' and '%s' in parents and trashed = false",
			driveFolderMIME,
			folder.Name,
			currentParentID,
		)

		res, err := cfg.service.Files.List().
			Q(query).
			Fields("files(id, name, parents)").
			PageSize(1).
			Do()
		if err != nil {
			return nil, fmt.Errorf("error querying folder '%s': %w", folder.Name, err)
		}

		if len(res.Files) == 0 {
			return nil, fmt.Errorf("%w: %s", ErrDirNotFound, folder.Name)
		}

		currentFolder = res.Files[0]
		currentParentID = currentFolder.Id
	}

	return currentFolder, nil
}
```
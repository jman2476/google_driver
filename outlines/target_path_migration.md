# Target Path Migration & File Renaming Checklist

This checklist outlines the locations and logic across the codebase requiring fixes or updates to complete the migration of the upload function signature from target directory to target path, supporting file renaming upon upload, handling destination directories without dots, and correctly processing files without extensions.

---

### 1. Upload Target Path Parsing & Disambiguation (`upload.go`)

- [X] `Lines 40-45` - Handle empty target path inputs so uploads with no target destination default to the root folder with the original filename instead of failing path validation when `ParsePath` returns an empty directory error.

- [ ] `Lines 40-56` - Support explicit directory paths indicated by a trailing slash (e.g., specifying a folder destination where the file keeps its original filename).

- [ ] `Lines 51-56`- Replace the period detection check on the final path segment (`strings.Contains(..., ".")`) with robust logic that can distinguish between a folder destination and a target filename.

- [ ] `Lines 51-56` - Handle destination folders that contain dots in their folder names so they are not misclassified as target filenames.

- [ ] `Lines 51-56` - Handle target filenames that have no file extension (such as documentation files, license files, or scripts) so they are recognized as file renames rather than directory names.

- [ ] `Lines 51-66` - Resolve target paths where only a new filename is provided with no directory path (uploading into the root directory under a new name).

- [ ] ` Lines 58-66` - Implement destination resolution against Drive to determine whether the path refers to an existing folder before falling back to extracting a parent folder and new filename.

- [ ] `Lines 58-76` - Set the extracted parent folder path and target filename variables appropriately so valid folder IDs are retrieved and passed into metadata creation.

---

### 2. File Metadata Construction (`upload.go`)

- [ ] `Lines 12-27` - Verify that `writeUploadMetadata` uses the new target name when specified (lines 22-24), while preserving the original local filename (lines 13-16) when no rename is requested.

- [ ] `Lines 13-16` - Ensure that MIME type detection performed during file validation `detect.go:58-72` is retained on the Drive metadata even when uploading or renaming to a filename without an extension.

- [ ] `Lines 18-20` - Ensure that parent folder ID assignment correctly targets the resolved parent folder ID (or root) regardless of whether a rename occurred.

---

### 3. Path Parsing Helpers (`directory.go`)

- [ ] `Lines 32-55` - Verify that `ParsePath` handles single-segment inputs, empty inputs, and paths with trailing slashes without returning unexpected empty directory errors `directory.go:50-52`.

- [ ] `Lines 57-110` - Ensure `findFolder` can search and verify parent folders independently from trailing file components in a target path.

---

### 4. Upload Function Testing (`upload_test.go`)

- [ ] `Lines 10-75` - Update `TestWriteUploadMetadata` table cases (lines 11-40) and function call (line 48) to pass explicit new filenames, verifying that the resulting metadata struct reflects the renamed value.

- [ ] `Lines 11-40` - Add `TestWriteUploadMetadata` unit test cases for files without extensions to confirm filename and MIME assignment.

- [ ] `Lines 89-107` - Update `TestUploadFile` integration test cases to test uploading into a target folder without renaming (retaining the original filename inside the folder).

- [ ] `Lines 89-107` - Add `TestUploadFile` integration test cases that provide a target path containing both a destination folder and a new filename to test folder placement and renaming simultaneously.

- [ ] `Lines 89-107` - Add `TestUploadFile` integration test cases that provide only a new filename to verify renaming directly within the root folder.

- [ ] `Lines 89-107` - Add `TestUploadFile` integration test cases for local files without extensions (such as license and test resources without extensions).

- [ ] `Lines 89-107` - Add `TestUploadFile` integration test cases that rename a file to a target name that does not have an extension.

---

### 5. Removal Integration Tests (`remove_test.go`)

- [ ] `Lines 55-67` - Update `uploadFile` calls in `TestRemoveFile` (line 57) so the target path argument passes the full destination path including the timestamped filename that `removeFile` subsequently attempts to delete (lines 62-63).

- [ ] `Lines 12-15` - Fix `genTimeStampFilename` (line 14) to ensure proper dot separation between the timestamp and the extension string.

- [ ] `Lines 33-53` - Add test cases in `TestRemoveFile` to verify removing files that do not have extensions.

- [ ] `Lines 126-132 - Update `TestRemoveFileErrors` cleanup routines so files uploaded during testing are removed from their actual destination folders.

---

### 6. CLI Command Registry, Dispatcher & Help (`command.go`, `main.go`, `command_help.go`, `parse.go`)

- [ ] `command.go Lines 23-27` - Verify the `mv` command in `commandRegistry` maps arguments to local source path and remote target path consistently with the dispatcher signature `command.go:8`.

- [ ] `command.go:25` & `command_help.go:31-43` - Update user-facing descriptions and help documentation for the `mv` command to clarify that the target parameter accepts either a destination folder path or a full destination path with a new filename.

- [ ] `parse.go:12-31` & `main.go:66-77` - Ensure argument parsing in `parseArgs` (lines 19-25) and invocation in `main` (line 67) provide an empty target string when optional target arguments are omitted by the user.

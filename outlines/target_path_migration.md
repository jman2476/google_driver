# Target Path Migration & File Renaming Checklist

This checklist outlines the locations and logic across the codebase requiring fixes or updates to complete the migration of the upload function signature from target directory to target path, supporting file renaming upon upload, handling destination directories without dots, and correctly processing files without extensions.

---

### 1. Upload Target Path Parsing & Disambiguation (`upload.go`)

--[ ] Handle empty target path inputs so uploads with no target destination default to the root folder with the original filename instead of failing path validation.
--[ ] Support explicit directory paths indicated by a trailing slash (e.g. specifying a folder destination where the file keeps its original filename).
--[ ] Replace the period detection check on the final path segment with robust logic that can distinguish between a folder destination and a target filename.
--[ ] Handle destination folders that contain dots in their folder names so they are not misclassified as target filenames.
--[ ] Handle target filenames that have no file extension (such as documentation files, license files, or scripts) so they are recognized as file renames rather than directory names.
--[ ] Resolve target paths where only a new filename is provided with no directory path (uploading into the root directory under a new name).
--[ ] Implement destination resolution against Drive to determine whether the path refers to an existing folder before falling back to extracting a parent folder and new filename.
--[ ] Set the extracted parent folder path and target filename variables appropriately so valid folder IDs are retrieved and passed into metadata creation.

---

### 2. File Metadata Construction (`upload.go`)

--[ ] Verify that metadata construction uses the new target name when specified, while preserving the original local filename when no rename is requested.
--[ ] Ensure that MIME type detection performed during file validation is retained on the Drive metadata even when uploading or renaming to a filename without an extension.
--[ ] Ensure that parent folder ID assignment correctly targets the resolved parent folder ID (or root) regardless of whether a rename occurred.

---

### 3. Path Parsing Helpers (`directory.go`)

--[ ] Verify that path splitting helpers handle single-segment inputs, empty inputs, and paths with trailing slashes without returning unexpected empty directory errors.
--[ ] Ensure folder lookup methods can search and verify parent folders independently from trailing file components in a target path.

---

### 4. Upload Function Testing (`upload_test.go`)

--[ ] Update metadata creation unit tests to include test cases that pass explicit new filenames, verifying that the resulting metadata struct reflects the renamed value.
--[ ] Add metadata creation unit test cases for files without extensions to confirm filename and MIME assignment.
--[ ] Update upload file integration test cases to test uploading into a target folder without renaming (retaining the original filename inside the folder).
--[ ] Add upload file integration test cases that provide a target path containing both a destination folder and a new filename to test folder placement and renaming simultaneously.
--[ ] Add upload file integration test cases that provide only a new filename to verify renaming directly within the root folder.
--[ ] Add upload file integration test cases for local files without extensions (such as license and test resources without extensions).
--[ ] Add upload file integration test cases that rename a file to a target name that does not have an extension.

---

### 5. Removal Integration Tests (`remove_test.go`)

--[ ] Update upload calls in file removal tests so the target path argument passes the full destination path including the timestamped filename that the subsequent removal step expects.
--[ ] Fix the timestamp generation helper to ensure proper dot separation between the timestamp and the extension string.
--[ ] Add test cases in removal tests to verify removing files that do not have extensions.
--[ ] Update removal error test cleanup routines so files uploaded during testing are removed from their actual destination folders.

---

### 6. CLI Command Registry, Dispatcher & Help (`command.go`, `main.go`, `command_help.go`)

--[ ] Verify the move/upload command registration maps arguments to local source path and remote target path consistently with the dispatcher.
--[ ] Update user-facing descriptions and help documentation for the move/upload command to clarify that the target parameter accepts either a destination folder path or a full destination path with a new filename.
--[ ] Ensure argument parsing in the main entry point provides an empty target string when optional target arguments are omitted by the user.


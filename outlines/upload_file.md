Yes, your understanding is spot on!

In the Google Drive v3 Go SDK, the client library uses a builder pattern to separate file metadata configuration from payload handling and request execution.

Here is how the methods fit together conceptually:

### 1. Creating a File (`Files.Create`)
* **Metadata Initialization:** You pass a pointer to a `drive.File` struct into `Files.Create()`. This struct only holds the file's metadata attributes (such as the display name, parent folder IDs, description, and metadata MIME type).
* **Attaching Content (`.Media()`):** Calling `Files.Create()` returns a builder object (`*FilesCreateCall`). You use the `.Media()` method on this call object to attach your actual file payload. The `.Media()` method accepts any standard reader stream (such as an opened local file or in-memory buffer). You can also pass optional parameters into `.Media()` to specify media content types or upload chunk options.
* **Executing the Request (`.Do()`):** Finally, you chain `.Do()` at the end of the call builder sequence. This performs the network request, sends the metadata and media payload to Google Drive, and returns the created file's remote metadata or an error.

### 2. Updating an Existing File (`Files.Update`)
* `Files.Update()` follows the exact same pattern. It accepts the target file's ID along with updated metadata and returns an update call builder object (`*FilesUpdateCall`).
* If you want to update the actual file content, you chain `.Media()` on that update call builder with your new reader stream, and then finish by calling `.Do()`.
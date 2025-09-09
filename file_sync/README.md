We are building a file sync application in Go. This involves several key components and challenges, such as:

1. Detecting changes in files and directories (both on the local system and remote server).
2. Efficiently transferring only the changed parts of files (to save bandwidth and time).
3. Handling conflicts when the same file is modified in two different places.
4. Ensuring reliability and fault tolerance (resume after interruption, etc.).
Here's a high-level plan:

### Algorithms and Patterns to Consider

1. **Change Detection**:
   - **Inotify/Watcher**: Use filesystem notification mechanisms (like `fsnotify` in Go) to monitor changes in real-time. Alternatively, periodic scanning (like rsync) can be used.
   - **Snapshot Comparison**: Maintain a snapshot (metadata of files: path, size, modification time, and hash) of the last synchronized state. Compare the current state with the snapshot to detect changes.
2. **Delta Encoding / Efficient Transfer**:
    - **Rsync Algorithm**: This algorithm efficiently transfers files by only sending the parts of the file that have changed. It breaks the file into fixed-size blocks and computes a rolling checksum for each block. Then, it compares the checksums to identify changed blocks. The `rsync` algorithm is well-known and there are Go implementations (like `github.com/gokrazy/rsync`).
   - **Content-Defined Chunking (CDC)** (like in Dropbox): Instead of fixed-size blocks, break files into variable-sized chunks using a sliding window (e.g., using Rabin fingerprinting). This way, insertions and deletions only affect a few chunks. Then, transfer only the chunks that are new.
3. **Conflict Resolution**:
   - **Last Write Wins (LWW)**: Simple but can lead to data loss.
   - **Version Vectors / Vector Clocks**: Track the history of changes to detect conflicts and allow manual resolution or more sophisticated automatic merging (for text files, we can use diff3, but for binary files, it's harder).
   - **Operational Transformation (OT) or Conflict-Free Replicated Data Type (CRDT)**: These are more advanced and typically used in collaborative editing. For a file sync app, they might be overkill, but if you plan to have real-time collaboration, they could be considered.
4. **Network Communication**:
   - **Client-Server or Peer-to-Peer?**: Depending on your architecture.
   - **Protocol**: You can use HTTP, WebSocket, or a custom protocol over TCP. For efficiency, consider binary protocols (like gRPC with protocol buffers) instead of JSON.
5. **Concurrency and Parallelism**:
   - **Goroutines and Channels**: Go's concurrency primitives are excellent for handling multiple file transfers concurrently and managing the filesystem watcher without blocking.
   - **Worker Pools**: Use a pool of workers to handle file transfers in parallel.
6. **Resumable Transfers**:
   - **Checkpointing**: For large files, store the state of the transfer (which blocks have been transferred) so that if interrupted, the transfer can resume from the last checkpoint.
7. **Security**:
   - **TLS**: Encrypt data in transit.
   - **Authentication**: Use tokens, OAuth, or public/private keys to authenticate clients.
8. **Storage**:
   - **Local Database**: To store snapshots and synchronization state, you can use an embedded database like SQLite or BoltDB.

### High-Level Design

1. **Client and Server**:
   - The client runs on the user's machine and monitors the local folder.
   - The server (if using a central server) holds the master copy and coordinates sync between multiple clients.
2. **Synchronization Steps**:
   a. **Indexing**: Build a snapshot of the current state (list of files with metadata and hashes).
   b. **Comparison**: Compare the current state with the last known state (from the server or from the last sync) to generate a list of changes.
   c. **Transfer**: Transfer the changes (new/changed/deleted files) in the most efficient way.
   d. **Conflict Detection and Resolution**: During comparison, if a file has been modified on both sides, flag a conflict and resolve (automatically or manually).
   e. **Update Snapshot**: After successful sync, update the snapshot to the new state.
3. **File Transfer**:
   - For new or changed files, break them into chunks and compute checksums (either fixed-size or variable-size).
   - Send the list of checksums to the server. The server responds with which chunks it already has (by comparing with its own checksums) and which it needs.
   - The client then sends only the needed chunks.
### Recommended Libraries in Go
- **Filesystem Watching**: `github.com/fsnotify/fsnotify`
- **Rsync-like Functionality**: `github.com/gokrazy/rsync` (or implement your own using rolling checksums)
- **CDC (Content-Defined Chunking)**: You can use `github.com/restic/chunker` for Rabin fingerprinting.
- **Database**: `github.com/boltdb/bolt` (now `go.etcd.io/bbolt`) or `github.com/mattn/go-sqlite3` for SQLite.
- **Networking**: Use gRPC (for structured RPC) or standard net/http for a REST-like API.
### Example Flow for a Client
1. **Initialization**:
   - Load last snapshot (if exists) from local DB.
   - Start a filesystem watcher to monitor changes.
2. **On Change or Periodic Sync**:
   - Build the current state (by scanning the directory tree, if using periodic sync, or by updating from watcher events and then building a snapshot to avoid missing events during processing).
   - Compare the current state with the last snapshot to create a change set (new, modified, deleted).
   - For each changed file, compute the chunk checksums (if using delta encoding) and check with the server to see which chunks are needed.
   - Transfer the missing chunks and reconstruct the file on the server side.
   - On the server, apply the changes and then broadcast to other clients (if any) that a change has occurred.
3. **Handling Deletes**:
   - Mark the file as deleted in the snapshot and send a delete command to the server. The server then propagates the delete to other clients.
### Conflict Resolution Strategy
- When a file is modified on both sides since the last sync, we have a conflict.
- One simple way: rename one of the files (e.g., `file.conflict-<timestamp>-<clientID>` and keep the other) and let the user decide.
- Alternatively, if the file is a text file, try to merge using a three-way merge (using the common ancestor). There are Go libraries for diffing and merging, like `github.com/sergi/go-diff/diffmatchpatch`.
### Summary of Patterns
- **Observer Pattern**: For filesystem events (using `fsnotify`).
- **Producer-Consumer**: Use a channel to queue changes detected by the watcher, and a pool of workers to process them.
- **Delta Encoding**: Use the rsync algorithm or CDC for efficient transfer.
- **State Machine**: Each file can be in a state (synced, syncing, conflicted, etc.) and the sync logic transitions these states.
### Example Code Structure
```
project/
├── client/
│   ├── main.go         # client entry point
│   ├── watcher.go      # filesystem watcher
│   ├── scanner.go      # periodic full scan (if needed)
│   ├── snapshot.go     # manage snapshot (current state and last synced state)
│   ├── sync.go         # synchronization logic (compare, transfer, etc.)
│   └── db.go           # database interaction (for snapshots)
├── server/
│   ├── main.go         # server entry point
│   ├── service.go      # gRPC or HTTP service
│   ├── storage.go      # handle file storage and chunk management
│   └── db.go           # server's database (if needed)
└── protocol/
    └── sync.proto      # if using gRPC, define the messages and services
```
### Conclusion
This is a complex project, so start with a minimal viable product (MVP). For example:
1. Implement a simple sync that transfers entire files (without chunking) and uses periodic scanning.
2. Then add filesystem watching.
3. Then add efficient transfer (rsync or CDC).
4. Then add conflict resolution.
Remember to write tests and consider edge cases (like file locks, permissions, etc.).
Good luck with your project!
To develop a robust file synchronization application in Go, you'll need a combination of efficient algorithms, concurrency patterns, and system design principles. Here's a structured approach:

# Build for Windows:
```bash
# Build normally (will include service functionality)
GOOS=windows GOARCH=amd64 go build -o file-sync.exe
```

# Service Management (Run as Administrator):
```cmd
# Install the service
file-sync.exe install

# Start the service
file-sync.exe start

# Stop the service
file-sync.exe stop

# Uninstall the service
file-sync.exe uninstall

# Run as console app (normal mode)
file-sync.exe

# Run without arguments - automatically detects if running as service or console
file-sync.exe
```

# Using sc command:
```cmd
# Check service status
sc query FileSyncService

# Start service
sc start FileSyncService

# Stop service  
sc stop FileSyncService
```

## Remember: After building, run as Administrator for service commands:

```cmd
# Run as Administrator in Command Prompt
file-sync.exe install
file-sync.exe start
```

## If you get execution policy error:

```powershell
# Run this first to allow script execution
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```
## Creating to the .exe file

```
PS <pwd>> .\scripts\build.ps1
>> 
Building Windows executable...
Build successful! Created file-sync.exe
Press Enter to continue: 
```
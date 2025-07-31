# 📝 Task Tracker

A command-line tool to manage daily tasks, store them in JSON, and retrieve them based on dates or status.

---

## ✅ Functionality

The application runs entirely from the **command line**. It accepts user actions and inputs as **flags/arguments**, and stores all tasks in a `tasks.json` file located in the current directory.

### 🔧 Supported Operations:
- **Add** a task  
- **Update** a task  
- **Delete** a task  
- **Mark** a task as _in progress_ or _done_  
- **List** all tasks  
- **List tasks by date**  

### 🗃️ Task Storage
- Tasks are indexed by ID in a JSON file.
- Each task includes metadata such as title, description, status, and timestamps.
- You can retrieve tasks based on a specific date.

---

## ⚙️ Constraints

- Uses **only the standard library** (no third-party dependencies).
- Uses **flags** (`-h`, `-add`, `-update`, etc.) for CLI usage.
- JSON file is **automatically created** if it doesn’t exist.
- Uses the **native file system module** for I/O operations.
- Implements **error handling** and **graceful fallbacks**.

---

## 🧪 To-Do: Write All Tests

> Use Go's `testing` package to implement the following:

- [ ] Test task creation (valid and invalid inputs)
- [ ] Test updating a task
- [ ] Test deleting a task
- [ ] Test marking task as done/in progress
- [ ] Test listing tasks (all and by date)
- [ ] Test JSON file creation and loading
- [ ] Test ID indexing and uniqueness
- [ ] Test edge cases (empty file, invalid JSON, corrupted task data)
- [ ] Test command-line flag parsing
- [ ] Test behavior when file is missing or read-only

---

> Made with 💻 in Go — no external dependencies.

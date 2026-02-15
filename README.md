# QuickTick (qt)

QuickTick is a terminal-based task manager built in Go, designed for speed and simplicity. It supports local data storage with SQLite and optional cloud synchronization via Supabase.

## Features

- **CLI Interface**: Efficient task management directly from the terminal.
- **Local Storage**: Utilizes SQLite for low-latency operations.
- **Multi-User Support**: Isolated databases for individual users.
- **Cloud Synchronization**: Optional two-way sync with Supabase (PostgreSQL).
- **Task Notes**: Attach markdown notes to tasks.

## Installation

### Prerequisites

- Go 1.22 or higher.

### Build via Source

```bash
git clone https://github.com/rafawastaken/quicktick.git
cd quicktick
go build -o qt.exe ./cmd/qt
```

Ensure `qt.exe` is in your PATH for global access.

## Usage

### Core Commands

```bash
# Add a new task
qt --add "Task description"

# List tasks
qt --show
qt --show --status todo
qt --show --status completed

# Complete a task
qt --done <ID>

# Edit a task
qt edit <ID> --title "New description"

# Delete a task
qt rm <ID>

# Open task notes (defaults to system editor)
qt --open <ID>
```

### Configuration and Sync

To enable cloud synchronization, configure the application with your Supabase credentials.

```bash
# Set credentials
qt config --url <SUPABASE_URL> --key <SUPABASE_KEY>

# Create an account
qt signup

# Login
qt login

# Synchronize tasks
qt --sync
```

## Database Setup

This application requires a Supabase PostgreSQL database for synchronization.

Please execute the SQL files located in the `migrations/` directory to set up the necessary tables and permissions.

1.  Create tables (schema definition).
2.  Enable Row Level Security (RLS) policies and user grants.

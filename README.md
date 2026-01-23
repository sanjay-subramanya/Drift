# Drift CLI

**Drift** is a lightweight, high-performance CLI tool designed to detect "code drift" between your local environment and upstream branches. It identifies stale files and categorizes changes by severity before you even attempt a `git pull`. It can compare staleness across:
- **Any Remote Branch**: Comparing your local changes (unstaged/staged) against any remote branch.
- **Forked Repositories**: Detecting local changes from the original upstream repository (not just your personal fork).


## Why Drift?

In fast-moving repositories, your local branch can become "stale" within hours. Standard Git commands tell you *if* you are behind, but **Drift** tells you **what** is different and **how critical** those changes are to your current work.

- ⚡ **Zero-Fetch Analysis**: Compare states without force-pulling objects.
- 🔍 **Severity Scoring**: Automatically categorizes drift as `CRITICAL`, `HIGH`, or `LOW`.
- 🤖 **CI/CD Ready**: Supports `--json` output for automated pipeline gates.
- 🔌 **Extensible**: Powering the [GitDrift VS Code Extension](https://marketplace.visualstudio.com/items?itemName=sanjay-subramanya.gitdrift).


## 🛡️ Conflict Prediction

One of Drift's most powerful features is its ability to **predict merge conflicts** before they happen. 

By cross-referencing upstream changes with your currently modified (unstaged/staged) files, Drift flags potential "Collision Zones." If a file you are currently editing has been modified on the server, Drift marks it as *[CRITICAL]*, allowing you to rebase or communicate with teammates before you deal with a messy manual merge.


## 🌳 Cross-Fork Detection

If you are working on a local clone of a forked open-source project, Drift can compare your local changes against the *original upstream repository's base branch* (e.g., `main` or `master`), not just your personal fork's upstream. This allows you to proactively identify stale files or potential overlaps with the main project before you even consider submitting a pull request, ensuring your contributions are always up-to-date and conflict-free.


## Installation

Run the following commands to install the `drift` CLI on your system:

- Windows (PowerShell):
```powershell
iwr https://raw.githubusercontent.com/sanjay-subramanya/drift/main/install.ps1 | iex
```

- MacOS / Linux (Shell):
```bash
curl -sSL https://raw.githubusercontent.com/sanjay-subramanya/drift/main/install.sh | sh
```

## Manual Installation

If you have **Go** installed on your system, you can build the binary manually. This is the recommended method for developers or those who want to verify the source code.

### 1. Build the Binary
Clone the repository and run the build command from the root directory:

```bash
# Clone the repository
git clone https://github.com/sanjay-subramanya/drift.git
cd drift
```

- Windows
```bash
go build -o drift.exe ./cmd/drift
```

- Mac/Linux
```bash
go build -o drift ./cmd/drift
```

### 2. Add to PATH
Add the executable file to your system's PATH to make it accessible from any directory.

- **Windows**: Add the directory containing `drift.exe` to your `PATH` environment variable.

- **MacOS / Linux**: Add the directory containing `drift` to your `PATH` in your shell configuration file (e.g., `~/.bashrc`, `~/.zshrc`).

### 3. Verify Installation
Run `drift --help` to confirm that the installation was successful.

## Usage
Once installed, you can use Drift in any Git repository. Simply navigate to your project directory and run:

```bash
drift
```
Drift will analyze the current branch with respect to the remote branch (default: `origin/main`) and report any detected code drift. If there are any conflicts, it will flag them with severity levels.
If you want to analyze a different branch, you can use the `--base flag:

```bash
drift --base branch_name
```

If you are working on a forked repository, you can use the `--upstream-url` flag to specify the original upstream repository's URL. This allows Drift to compare your local changes against the base branch of the upstream repository.

 ```bash
drift --upstream-url https://github.com/original-owner/original-repo.git --base branch_name
```


For storing the output in a JSON file, you can use the `--json` flag and specify the path (default: `./drift.json`):

```bash
drift --json --path=drift_output.json
```

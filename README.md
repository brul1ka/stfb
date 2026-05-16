# STFB - Sort these files, baby! CLI utility for Windows and Linux to sort a bunch of files into generated directories.
## Features:
* **-undo** - if the files were sorted in a different place than planned, the action can be undone.
* **-dry-run** - you can recheck what will be sorted where.
* **Configuration** - you can configure config file for, e.g., add custom exts rules. The app looks for `config.yaml` alongside the file being sorted
## Installation
*You need Go 1.21 or higher installed on your system.*
```bash
git clone https://github.com/brul1ka/stfb.git
cd stfb
go build -o stfb .
```
Optional: Install system-wide (Linux/macOS):
```bash
sudo mv stfb /usr/local/bin/
```
## Usage
Basic sorting (scans current directory and moves files into organized folders):
```bash
./stfb
```
Sort files in a specific directory:
```bash
./stfb -path /home/user/Downloads
```
Check what will happen without actually moving anything:
```bash
./stfb -path /home/user/Downloads -dry-run
```
Undo last sorting:
```bash
./stfb -undo
```
### How it works:
**Before:**
```text
Downloads/
├── report.pdf
├── photo.jpg
├── script.py
└── todo.md
```
**After running** `stfb`:
```text
Downloads/
├── Office/
│   └── report.pdf
├── Pictures/
│   └── photo.jpg
├── Code/
│   └── script.py
└── Documents/
    └── todo.md
```
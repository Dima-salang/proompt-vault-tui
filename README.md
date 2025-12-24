# Proompt Vault

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Proompt Vault is a side project of mine because I found that I was always losing my system prompts and I had to recreate them every time on the fly. I wanted to have a simple interface where I can see all of my prompts and efficiently filter through so I can copy the one I need.

## Why?

I wanted something that:
- **Is fast**: Fuzzy search is a must. I don't want to type exact titles.
- **Supports Vim keys**: I'm a Vim user and I find it more efficient to use Vim keys to navigate the interface.
- **Stays local**: No accounts, no cloud sync, just a simple `prompts.db` file on your machine.

## How to use it

### Installation

If you have Go installed, just clone and build.

```bash
git clone https://github.com/Dima-salang/proompt-vault-tui.git
cd proompt-vault-tui
go build -o pvt .
```


(Optional) Toss it in your path so you can run it from anywhere:
```bash
sudo mv pvt /usr/local/bin/
```

Or you can do go get it:
```bash
go install "github.com/Dima-salang/proompt-vault-tui@latest"
```

This will install it to your `$GOPATH/bin` directory. 
Since the bin executable name is quite long, you can rename it with pvt or any name you like for easy access.
```
mv $GOBIN/proompt-vault-tui $GOBIN/pvt
```

### Running it

Just type:
```bash
pvt
```

### Controls

The interface is pretty intuitive and supports Vim keys for navigating up and down.

**In the List:**
- `↑` / `↓` or **Mouse Wheel**: Scroll through your collection.
- `/`: Start typing to fuzzy search.
- `Enter`: **Copy the selected prompt**. This is the main action.
- `a`: Add a new one.
- `e`: Edit the one you're hovering over.
- `d`: Delete it (with a confirmation check, don't worry).

**In the Editor:**
- `Tab` / `Shift+Tab`: Move between fields.
- `Enter` (on the Submit button): Save it.
- `Esc`: Cancel and go back.


## Under the hood

This is a pure Go project. I used the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework because it's awesome for building TUIs. Styling is handled by [Lip Gloss](https://github.com/charmbracelet/lipgloss), and the data lives in [BoltDB](https://github.com/boltdb/bolt) (a solid key/value store).

## 🤝 Contributing

If you want to add features (like export/import or maybe syntax highlighting), feel free to fork it and open a PR. I'm open to ideas!

---

# paq

跨平台 CLI 工具，統一查詢多個套件管理器已安裝的套件。

[English](README.md)

## 為什麼需要 paq？

開發者通常安裝了多個套件管理器，每個都有不同的指令：

```bash
brew list          # Homebrew
npm list -g        # npm
snap list          # Snap
flatpak list       # Flatpak
```

**paq** 提供一致的統一介面：

```bash
paq brew
paq npm
paq snap
paq flatpak
```

## 功能特色

- **單一二進位檔** — 跨平台（macOS、Linux、Windows）
- **零設定** — 自動偵測可用的套件管理器
- **雙格式輸出** — 表格（預設）與 JSON（`--json`）
- **可擴展** — 新增 Provider 不需修改核心程式碼

## 安裝

### Homebrew

```bash
brew install cluion/tap/paq
```

### Go install

```bash
go install github.com/cluion/paq@latest
```

### 從原始碼編譯

```bash
git clone https://github.com/cluion/paq.git
cd paq
make build
sudo make install   # 安裝至 /usr/local/bin/paq
```

如果不想使用 `sudo`，可以在 shell 設定檔（`~/.zshrc` 或 `~/.bashrc`）加入 alias：

```bash
alias paq="/path/to/paq/bin/paq"
```

### 下載二進位檔

從 [Releases](https://github.com/cluion/paq/releases) 下載最新版本。

## 更新

**Go：**

```bash
go install github.com/cluion/paq@latest
```

**Homebrew：**

```bash
brew upgrade cluion/tap/paq
```

**從原始碼：**

```bash
cd paq
git pull
make build
sudo make install   # 使用 alias 的話，重新編譯即可立即生效
```

## 發佈流程

維護者推送 tag 即可自動觸發 GitHub Actions 建置、發佈 Release、更新 Homebrew formula：

```bash
git tag v1.0.0
git push origin v1.0.0
```

## 使用方式

### 列出偵測到的套件源

```bash
paq list
```

```
Source  Display Name  Status
brew    Homebrew      available
npm     npm           available
snap    Snap          not found
flatpak Flatpak       not found
Total: 4
```

### 查詢特定套件源

```bash
paq brew
```

```
Name        Version   Description
git         2.40.0    Distributed version control system
node        20.1.0    Platform built on V8
curl        8.1.0
Total: 3
```

### JSON 格式輸出

```bash
paq brew --json
```

```json
{
  "source": "brew",
  "count": 3,
  "packages": [
    { "name": "git", "version": "2.40.0", "desc": "Distributed version control system" },
    { "name": "node", "version": "20.1.0", "desc": "Platform built on V8" },
    { "name": "curl", "version": "8.1.0" }
  ]
}
```

### 版本資訊

```bash
paq version
```

```
paq dev (commit: abc1234, built: 2026-04-06)
```

## 支援的套件管理器

| 套件管理器 | 指令 | 平台 |
|-----------|------|------|
| Homebrew | `paq brew` | macOS |
| npm | `paq npm` | macOS、Linux、Windows |
| Snap | `paq snap` | Linux |
| Flatpak | `paq flatpak` | Linux |

## 開發

### 環境需求

- Go 1.22+
- golangci-lint

### 編譯

```bash
make build
```

### 測試

```bash
make test
```

### Lint

```bash
make lint
```

## 授權

[MIT](LICENSE)

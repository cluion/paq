# 更新日誌

本專案所有值得注意的變更都會記錄在此文件中。

格式基於 [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)，
本專案遵循 [語意化版本](https://semver.org/spec/v2.0.0.html)。

## [1.2.1] - 2026-04-07

### 修正

- `paq upgrade` 所有安裝方式統一加入版本比對（Homebrew、go install、自我更新）
- 修正 version 命令的 errcheck lint 錯誤
- CI Go 版本更新至 1.26 以配合 go.mod

### 新增

- CodeQL 安全性分析
- govulncheck 漏洞掃描
- CI 覆蓋率門檻檢查（80%）

## [1.2.0] - 2026-04-07

### 修正

- `paq upgrade` 自我更新前先比對版本再下載
- 自我更新前檢查寫入權限，避免無聲失敗

## [1.1.0] - 2026-04-07

### 新增

- `paq upgrade` 命令，自動偵測安裝方式並升級（Homebrew、go install、本地安裝、直接下載）
- `paq upgrade --check` 旗標，僅檢查新版本不執行升級
- 自我更新功能，從 GitHub Releases 下載 binary 直接替換（`minio/selfupdate`）
- Makefile `uninstall` 目標
- README 中英版新增解除安裝說明
- `debug.ReadBuildInfo()` fallback 版本顯示機制

### 變更

- 修正 `go install` 路徑：`github.com/cluion/paq@latest` → `github.com/cluion/paq/cmd/paq@latest`
- 版本顯示改用 `debug.ReadBuildInfo()` 作為 fallback，修正 `go install` 版本顯示問題

## [1.0.0] - 2026-04-06

### 新增

- paq 首次發布
- 統一 CLI 查詢多個套件管理器的已安裝套件
- Homebrew provider（`paq brew`）
- npm provider（`paq npm`）
- Snap provider（`paq snap`）
- Flatpak provider（`paq flatpak`）
- `paq list` 命令，顯示偵測到的套件源
- `paq version` 命令，含建置資訊注入
- 表格輸出格式（預設）
- JSON 輸出格式（`--json` 旗標）
- 自動偵測可用的套件管理器
- 跨平台支援（macOS、Linux、Windows）
- 測試覆蓋率 >= 80%
- Makefile（build、test、lint、clean、install targets）
- GoReleaser 自動化發佈設定
- GitHub Actions CI/CD pipeline

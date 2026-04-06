# 更新日誌

本專案所有值得注意的變更都會記錄在此文件中。

格式基於 [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)，
本專案遵循 [語意化版本](https://semver.org/spec/v2.0.0.html)。

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

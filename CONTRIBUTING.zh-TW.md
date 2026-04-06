# 貢獻指南

感謝你對 paq 的貢獻興趣！本文檔提供貢獻的相關規範。

[English](CONTRIBUTING.md)

## 如何貢獻

### 回報 Bug

開一個 [issue](https://github.com/cluion/paq/issues/new)，附上：

- 作業系統與架構
- paq 版本（`paq version` 輸出）
- 重現步驟
- 預期行為與實際行為

### 功能建議

開一個 [issue](https://github.com/cluion/paq/issues/new)，標記為 `enhancement`。

### 提交程式碼

1. Fork 此倉庫
2. 建立功能分支（`git checkout -b feat/my-feature`）
3. 進行修改
4. 確保所有檢查通過：
   ```bash
   make lint
   make test
   ```
5. 使用 [conventional commits](https://www.conventionalcommits.org/) 格式提交：
   ```
   feat: add cargo provider
   fix: handle empty npm output
   docs: update README
   ```
6. Push 並開啟 Pull Request

## 開發環境設定

### 環境需求

- Go 1.22+
- [golangci-lint](https://golangci-lint.run/usage/install/)

### 編譯與測試

```bash
make build    # 編譯至 bin/paq
make test     # 執行測試（含 race detector 與覆蓋率）
make lint     # 執行 golangci-lint
make clean    # 清除建置產物
```

### 專案結構

```
cmd/paq/main.go              # 程式入口
internal/
  cli/                        # CLI 命令（root, list, query, version）
  provider/                   # Provider 介面 + 實作
  output/                     # 格式化輸出（table, JSON）
Makefile
.goreleaser.yml
```

### 新增 Provider

1. 建立 `internal/provider/<name>.go`
2. 實作 `Provider` 介面：
   ```go
   type Provider interface {
       Name() string
       DisplayName() string
       Detect() bool
       List() ([]Package, error)
   }
   ```
3. 在 `init()` 中註冊：
   ```go
   func init() {
       Register(&MyProvider{runner: defaultRunner})
   }
   ```
4. 使用 fake `CommandRunner` 撰寫測試
5. 更新 README 與 CHANGELOG

### 程式碼風格

- 遵循 `gofmt` 與 `goimports`
- 錯誤訊息使用小寫，不加句號
- 介面保持小巧（1-3 個方法）
- 函式不超過 50 行，檔案不超過 800 行

## 授權

提交貢獻即表示你同意你的貢獻將以 [MIT License](LICENSE) 授權。

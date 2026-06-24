# 吃瓜資料庫 前端

Vite + Vue 3 實作的 V1 前端，對接後端兩支 API：

- `POST /api/gua` 存瓜
- `POST /api/gua/list` 查瓜

兩個表單皆有 `title`、`people`、`content` 三個欄位，**皆可留空**。
其中 `people` 是文字陣列，在前端以「逗號」分隔輸入（中英逗號皆可），送出時自動切成陣列。

## 啟動

需先啟動後端（預設 `localhost:8080`）。

```bash
cd frontend
npm install
npm run dev      # http://localhost:5173
```

開發時 `/api` 會透過 Vite proxy 轉到 `http://localhost:8080`，無跨網域問題。

## 打包

```bash
npm run build    # 輸出 dist/
npm run preview  # 預覽打包結果
```
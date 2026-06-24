## 吃瓜資料庫

## 功能簡介
可以輸入瓜的內容，為瓜取名，並且列出瓜參與的人物，可以用關鍵查詢某個瓜，用人名去查他參與的瓜。     
未來可以加入共享瓜資料庫，一個人可以擁有很多個瓜資料庫，一個瓜資料庫可以被很多人擁有，就像是群組一樣。

## v1
建立完整前端後端。
可以存瓜可以查瓜。
不做共享資料庫。
不碰分詞器 (V2再做)
不做資料修改。

## schema V1
```sql
CREATE TABLE IF NOT EXISTS guas(
    gua_id BIGSERIAL PRIMARY KEY,
    title TEXT,
    people TEXT[],
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```
## API V1

POST /api/gua 新增瓜
接收資料

```json
{
  "data": {
    "title": "",
    "people": [""],
    "content" : ""
  }
}
```
回傳資料

```json
{
  "data": {
    "guaId": ""
  }
}
```
錯誤時回傳

```json
{
  "error": ""
}
```
POST /api/gua/list 查詢瓜
titile,people,content都是關鍵字有匹配到就回傳
```json
{
  "data": {
    "title": "",
    "people": [""],
    "content" : ""
  }
}
```
回傳資料

```json
{
  "data": [
    {
      "title": "",
      "people": [""],
      "content": ""
    }
  ]
}
```
錯誤時回傳

```json
{
  "error": ""
}
```

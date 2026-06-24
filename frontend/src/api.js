// 後端 API 封裝。V1 只有兩支：存瓜、查瓜。
// 兩支都是 POST，request body 用 { data: {...} } 包起來，
// 成功回 { data: ... }，失敗回 { error: "..." }。

// 把逗號分隔的字串切成人名陣列，去空白並過濾空字串。
export function parsePeople(text) {
  if (!text) return []
  return text
    .split(/[,，]/)
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
}

async function request(path, payload) {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ data: payload })
  })
  let body = null
  try {
    body = await res.json()
  } catch (e) {
    throw new Error(`伺服器回應不是有效 JSON（HTTP ${res.status}）`)
  }
  if (body && body.error) {
    throw new Error(body.error)
  }
  if (!res.ok) {
    throw new Error(body && body.error ? body.error : `HTTP ${res.status}`)
  }
  return body.data
}

// 存瓜。回傳新瓜的 guaId。
export function addGua({ title, people, content }) {
  return request('/api/gua', { title, people, content })
}

// 查瓜。回傳符合的瓜清單。
export function listGua({ title, people, content }) {
  return request('/api/gua/list', { title, people, content })
}
<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { addGua, listGua, parsePeople } from './api.js'

const tab = ref('save') // 'save' | 'query'

// 存瓜表單
const saveForm = ref({ title: '', people: '', content: '' })
const saveState = ref({ status: 'idle', message: '' }) // status: idle | loading | ok | error

// 查瓜表單
const queryForm = ref({ title: '', people: '', content: '' })
const queryState = ref({ status: 'idle', message: '' })
const results = ref([])

// 預覽截斷長度與彈窗
const PREVIEW_LIMIT = 120
const selected = ref(null) // 目前展開的瓜，null 表示關閉

function preview(content) {
  if (!content) return ''
  return content.length > PREVIEW_LIMIT
    ? content.slice(0, PREVIEW_LIMIT) + '...'
    : content
}

function openDetail(g) {
  selected.value = g
}

function closeDetail() {
  selected.value = null
}

function onKeydown(e) {
  if (e.key === 'Escape' && selected.value) closeDetail()
}
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

async function onSave() {
  const payload = {
    title: saveForm.value.title.trim(),
    people: parsePeople(saveForm.value.people),
    content: saveForm.value.content.trim()
  }
  saveState.value = { status: 'loading', message: '儲存中...' }
  try {
    const data = await addGua(payload)
    const id = data?.guaId ?? data?.gua_id
    saveState.value = { status: 'ok', message: `已儲存，瓜編號：${id}` }
    saveForm.value = { title: '', people: '', content: '' }
  } catch (e) {
    saveState.value = { status: 'error', message: e.message || '儲存失敗' }
  }
}

async function onQuery() {
  const payload = {
    title: queryForm.value.title.trim(),
    people: parsePeople(queryForm.value.people),
    content: queryForm.value.content.trim()
  }
  queryState.value = { status: 'loading', message: '查詢中...' }
  results.value = []
  try {
    const data = await listGua(payload)
    results.value = Array.isArray(data) ? data : []
    queryState.value = {
      status: 'ok',
      message: `找到 ${results.value.length} 筆`
    }
  } catch (e) {
    queryState.value = { status: 'error', message: e.message || '查詢失敗' }
  }
}
</script>

<template>
  <div class="page">
   <div class="wrap">
    <header class="header">
      <h1>吃瓜資料庫</h1>
      <p class="subtitle">存瓜、查瓜。三個欄位皆可留空。</p>
    </header>

    <nav class="tabs">
      <button
        :class="['tab', { active: tab === 'save' }]"
        @click="tab = 'save'"
      >
        存瓜
      </button>
      <button
        :class="['tab', { active: tab === 'query' }]"
        @click="tab = 'query'"
      >
        查瓜
      </button>
    </nav>

    <!-- 存瓜 -->
    <section v-if="tab === 'save'" class="panel">
      <h2>新增一筆瓜</h2>
      <form @submit.prevent="onSave">
        <label class="field">
          <span class="label">標題 title</span>
          <input v-model="saveForm.title" type="text" placeholder="可為空" />
        </label>
        <label class="field">
          <span class="label">人物 people</span>
          <input
            v-model="saveForm.people"
            type="text"
            placeholder="用逗號分隔多人，可為空。例：小明, 小華"
          />
        </label>
        <label class="field">
          <span class="label">內容 content</span>
          <textarea
            v-model="saveForm.content"
            rows="6"
            placeholder="可為空"
          ></textarea>
        </label>
        <div class="actions">
          <button type="submit" class="btn" :disabled="saveState.status === 'loading'">
            {{ saveState.status === 'loading' ? '儲存中...' : '存瓜' }}
          </button>
        </div>
      </form>
      <p
        v-if="saveState.message"
        :class="['msg', saveState.status === 'error' ? 'msg-error' : 'msg-ok']"
      >
        {{ saveState.message }}
      </p>
    </section>

    <!-- 查瓜 -->
    <section v-else class="panel">
      <h2>查詢瓜</h2>
      <p class="hint">任一欄位有匹配即回傳；留空表示不限制該欄位。</p>
      <form @submit.prevent="onQuery">
        <label class="field">
          <span class="label">標題 title</span>
          <input v-model="queryForm.title" type="text" placeholder="可為空" />
        </label>
        <label class="field">
          <span class="label">人物 people</span>
          <input
            v-model="queryForm.people"
            type="text"
            placeholder="用逗號分隔多人，可為空。例：小明"
          />
        </label>
        <label class="field">
          <span class="label">內容 content</span>
          <textarea
            v-model="queryForm.content"
            rows="3"
            placeholder="可為空"
          ></textarea>
        </label>
        <div class="actions">
          <button type="submit" class="btn" :disabled="queryState.status === 'loading'">
            {{ queryState.status === 'loading' ? '查詢中...' : '查瓜' }}
          </button>
        </div>
      </form>

      <p
        v-if="queryState.message"
        :class="['msg', queryState.status === 'error' ? 'msg-error' : 'msg-ok']"
      >
        {{ queryState.message }}
      </p>

      <ul v-if="results.length" class="results">
        <li v-for="(g, i) in results" :key="i" class="result">
          <div class="result-title">{{ g.title || '(無標題)' }}</div>
          <div v-if="g.people && g.people.length" class="result-people">
            人物：{{ g.people.join('、') }}
          </div>
          <pre v-if="g.content" class="result-preview">{{ preview(g.content) }}</pre>
          <button type="button" class="detail-btn" @click="openDetail(g)">
            查看更多
          </button>
        </li>
      </ul>
      <p v-else-if="queryState.status === 'ok'" class="empty">沒有符合的瓜</p>
    </section>
   </div>

    <!-- 完整內容彈窗 -->
    <div v-if="selected" class="overlay" @click.self="closeDetail">
      <div class="modal">
        <div class="modal-head">
          <h3>{{ selected.title || '(無標題)' }}</h3>
          <button type="button" class="modal-close" aria-label="關閉" @click="closeDetail"></button>
        </div>
        <div v-if="selected.people && selected.people.length" class="modal-people">
          人物：{{ selected.people.join('、') }}
        </div>
        <pre class="modal-content">{{ selected.content || '(無內容)' }}</pre>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page {
  padding: 24px 16px 48px;
  height: 100%;
  overflow-y: auto;
  overscroll-behavior-y: contain;
  background: var(--bg);
}
.wrap {
  max-width: 640px;
  margin: 0 auto;
}

.header h1 {
  margin: 0 0 4px;
  font-size: 22px;
}
.subtitle {
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.tabs {
  display: flex;
  gap: 8px;
  margin: 20px 0 16px;
}
.tab {
  padding: 8px 18px;
  border: 1px solid var(--border);
  background: var(--panel);
  color: var(--text);
  border-radius: 6px;
}
.tab.active {
  border-color: var(--accent);
  color: var(--accent);
  font-weight: 600;
}

.panel {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 20px;
}
.panel h2 {
  margin: 0 0 12px;
  font-size: 16px;
}
.hint {
  margin: 0 0 14px;
  color: var(--muted);
  font-size: 13px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 14px;
}
.label {
  font-size: 13px;
  color: var(--muted);
}
.field input,
.field textarea {
  width: 100%;
  padding: 9px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: #fff;
  color: var(--text);
  resize: vertical;
}
.field input:focus,
.field textarea:focus {
  outline: none;
  border-color: var(--accent);
}

.actions {
  margin-top: 4px;
}
.btn {
  padding: 9px 22px;
  border: 1px solid var(--accent);
  background: var(--accent);
  color: #fff;
  border-radius: 6px;
  font-weight: 600;
}
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.msg {
  margin: 14px 0 0;
  font-size: 14px;
}
.msg-ok {
  color: var(--accent);
}
.msg-error {
  color: var(--danger);
}

.results {
  list-style: none;
  margin: 20px 0 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.result {
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 12px 14px;
  background: #fafafa;
}
.result-title {
  font-weight: 600;
  margin-bottom: 6px;
}
.result-people {
  font-size: 13px;
  color: var(--muted);
  margin-bottom: 6px;
}
.result-content {
  margin: 0;
  font-family: inherit;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14px;
}
.result-preview {
  margin: 0 0 8px;
  font-family: inherit;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 14px;
  color: var(--text);
}
.detail-btn {
  border: 1px solid var(--border);
  background: #fff;
  color: var(--text);
  border-radius: 6px;
  padding: 5px 12px;
  font-size: 13px;
}
.detail-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  z-index: 50;
}
.modal {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 10px;
  max-width: 640px;
  width: 100%;
  max-height: 80vh;
  overflow-y: auto;
  overscroll-behavior-y: contain;
  padding: 20px;
}
.modal-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}
.modal-head h3 {
  margin: 0;
  font-size: 17px;
  word-break: break-word;
}
.modal-close {
  position: relative;
  flex-shrink: 0;
  border: 1px solid var(--border);
  background: #fff;
  color: var(--text);
  border-radius: 6px;
  width: 30px;
  height: 30px;
  padding: 0;
}
.modal-close::before,
.modal-close::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 14px;
  height: 2px;
  background: currentColor;
  border-radius: 1px;
}
.modal-close::before {
  transform: translate(-50%, -50%) rotate(45deg);
}
.modal-close::after {
  transform: translate(-50%, -50%) rotate(-45deg);
}
.modal-close:hover {
  border-color: var(--danger);
  color: var(--danger);
}
.modal-people {
  color: var(--muted);
  font-size: 13px;
  margin-bottom: 12px;
}
.modal-content {
  margin: 0;
  font-family: inherit;
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 15px;
  line-height: 1.6;
}
.empty {
  margin: 20px 0 0;
  color: var(--muted);
  font-size: 14px;
}
</style>
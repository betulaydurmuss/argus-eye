import {
    RunGoogleDork,
    GetDorks,
    GetCategories,
    GetHistory,
    ClearHistory,
    AddDork,
    DeleteDork,
    GetStats,
} from './wailsjs/go/main/App.js';

let allDorks = [];
let filteredDorks = [];
let activeCategory = '';
let searchQuery = '';
let activeTab = 'library';

const dorkGrid       = document.getElementById('dork-grid');
const targetInput    = document.getElementById('target');
const searchInput    = document.getElementById('search-input');
const categoryFilter = document.getElementById('category-filter');
const categoryTabs   = document.getElementById('category-tabs');
const dorkCountText  = document.getElementById('dork-count-text');
const emptyState     = document.getElementById('empty-state');
const toast          = document.getElementById('toast');

const statTotal   = document.getElementById('stat-total');
const statHistory = document.getElementById('stat-history');
const statCats    = document.getElementById('stat-cats');

const historyList      = document.getElementById('history-list');
const historyEmpty     = document.getElementById('history-empty');
const historyCountText = document.getElementById('history-count-text');
const clearHistoryBtn  = document.getElementById('clear-history-btn');

const newTitle      = document.getElementById('new-title');
const newQuery      = document.getElementById('new-query');
const newCategory   = document.getElementById('new-category');
const newTags       = document.getElementById('new-tags');
const previewText   = document.getElementById('preview-text');
const addDorkBtn    = document.getElementById('add-dork-btn');
const clearFormBtn  = document.getElementById('clear-form-btn');
const categoryList  = document.getElementById('category-list');

async function init() {
    await loadDorks();
    await refreshStats();
    setupEventListeners();
}

async function loadDorks() {
    try {
        allDorks = await GetDorks();
        const cats = await GetCategories();
        buildCategoryUI(cats);
        applyFilters();
    } catch (e) {
        showToast('Dorklar yüklenirken hata oluştu', 'error');
    }
}

async function refreshStats() {
    try {
        const stats = await GetStats();
        statTotal.textContent   = stats.totalDorks      ?? 0;
        statHistory.textContent = stats.totalHistory    ?? 0;
        statCats.textContent    = stats.totalCategories ?? 0;
    } catch (_) {}
}

function buildCategoryUI(cats) {
    categoryFilter.innerHTML = '<option value="">Tüm Kategoriler</option>';
    cats.forEach(cat => {
        const opt = document.createElement('option');
        opt.value = cat;
        opt.textContent = cat;
        categoryFilter.appendChild(opt);
    });

    categoryTabs.innerHTML = '';
    const allBtn = createCatTab('Tümü', '');
    categoryTabs.appendChild(allBtn);
    cats.forEach(cat => {
        categoryTabs.appendChild(createCatTab(cat, cat));
    });

    categoryList.innerHTML = '';
    cats.forEach(cat => {
        const opt = document.createElement('option');
        opt.value = cat;
        categoryList.appendChild(opt);
    });
}

function createCatTab(label, value) {
    const btn = document.createElement('button');
    btn.className = 'cat-tab' + (value === activeCategory ? ' active' : '');
    btn.textContent = label;
    btn.dataset.cat = value;
    btn.addEventListener('click', () => {
        activeCategory = value;
        document.querySelectorAll('.cat-tab').forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        categoryFilter.value = value;
        applyFilters();
    });
    return btn;
}

function applyFilters() {
    const q = searchQuery.toLowerCase();
    filteredDorks = allDorks.filter(d => {
        const matchCat   = !activeCategory || d.category === activeCategory;
        const matchSearch = !q ||
            d.title.toLowerCase().includes(q) ||
            d.query.toLowerCase().includes(q) ||
            d.tags.toLowerCase().includes(q);
        return matchCat && matchSearch;
    });

    dorkCountText.textContent = `${filteredDorks.length} / ${allDorks.length} sorgu gösteriliyor`;
    renderDorks();
}

function renderDorks() {
    dorkGrid.innerHTML = '';

    if (filteredDorks.length === 0) {
        emptyState.classList.remove('hidden');
        return;
    }
    emptyState.classList.add('hidden');

    filteredDorks.forEach(dork => {
        const card = document.createElement('div');
        card.className = 'card';
        card.dataset.id = dork.id;

        const isCustom = dork.id.startsWith('custom_');
        const tags = dork.tags
            ? dork.tags.split(',').map(t => `<span class="tag">${t.trim()}</span>`).join('')
            : '';

        card.innerHTML = `
            <div class="card-header">
                <span class="card-category">${escapeHtml(dork.category)}</span>
                ${isCustom ? '<span class="badge-custom">Özel</span>' : ''}
            </div>
            <h4 class="card-title">${escapeHtml(dork.title)}</h4>
            <code class="card-query">${escapeHtml(dork.query)}</code>
            ${tags ? `<div class="card-tags">${tags}</div>` : ''}
            <div class="card-actions">
                <button class="btn-run" data-id="${dork.id}" data-title="${escapeHtml(dork.title)}" data-query="${escapeHtml(dork.query)}">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    Sorgula
                </button>
                <button class="btn-copy" data-query="${escapeHtml(dork.query)}" title="Kopyala">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>
                </button>
                ${isCustom ? `<button class="btn-delete" data-id="${dork.id}" title="Sil">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14H6L5 6"/></svg>
                </button>` : ''}
            </div>
        `;
        dorkGrid.appendChild(card);
    });

    dorkGrid.querySelectorAll('.btn-run').forEach(btn => {
        btn.addEventListener('click', () => {
            executeDork(btn.dataset.title, btn.dataset.query);
        });
    });
    dorkGrid.querySelectorAll('.btn-copy').forEach(btn => {
        btn.addEventListener('click', () => {
            copyToClipboard(btn.dataset.query);
        });
    });
    dorkGrid.querySelectorAll('.btn-delete').forEach(btn => {
        btn.addEventListener('click', () => {
            deleteDork(btn.dataset.id);
        });
    });
}

async function executeDork(title, query) {
    const domain = targetInput.value.trim();
    try {
        const result = await RunGoogleDork(domain, title, query);
        if (result.success) {
            showToast(`"${title}" sorgusu tarayıcıda açıldı`, 'success');
            await refreshStats();
        } else {
            showToast(result.message, 'error');
        }
    } catch (e) {
        showToast('Sorgu çalıştırılırken hata oluştu', 'error');
    }
}

async function deleteDork(id) {
    try {
        const ok = await DeleteDork(id);
        if (ok) {
            showToast('Dork silindi', 'info');
            await loadDorks();
            await refreshStats();
        }
    } catch (e) {
        showToast('Silme işlemi başarısız', 'error');
    }
}

function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showToast('Sorgu panoya kopyalandı', 'info');
    }).catch(() => {
        showToast('Kopyalama başarısız', 'error');
    });
}

async function loadHistory() {
    try {
        const history = await GetHistory();
        historyCountText.textContent = `${history.length} arama kaydı`;

        if (history.length === 0) {
            historyList.innerHTML = '';
            historyEmpty.classList.remove('hidden');
            return;
        }
        historyEmpty.classList.add('hidden');
        historyList.innerHTML = '';

        history.forEach(entry => {
            const item = document.createElement('div');
            item.className = 'history-item';
            item.innerHTML = `
                <div class="history-meta">
                    <span class="history-time">${escapeHtml(entry.timestamp)}</span>
                    ${entry.domain ? `<span class="history-domain">${escapeHtml(entry.domain)}</span>` : '<span class="history-domain no-domain">Domain yok</span>'}
                </div>
                <div class="history-title">${escapeHtml(entry.dorkTitle)}</div>
                <code class="history-query">${escapeHtml(entry.fullQuery)}</code>
                <button class="btn-rerun" data-title="${escapeHtml(entry.dorkTitle)}" data-query="${escapeHtml(entry.query)}" data-domain="${escapeHtml(entry.domain)}">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/></svg>
                    Tekrar Çalıştır
                </button>
            `;
            historyList.appendChild(item);
        });

        historyList.querySelectorAll('.btn-rerun').forEach(btn => {
            btn.addEventListener('click', () => {
                targetInput.value = btn.dataset.domain;
                executeDork(btn.dataset.title, btn.dataset.query);
            });
        });
    } catch (e) {
        showToast('Geçmiş yüklenirken hata oluştu', 'error');
    }
}

function updatePreview() {
    const domain = targetInput.value.trim() || 'example.com';
    const query  = newQuery.value.trim() || '<sorgu buraya gelecek>';
    previewText.textContent = `site:${domain} ${query}`;
}

async function submitAddDork() {
    const title    = newTitle.value.trim();
    const query    = newQuery.value.trim();
    const category = newCategory.value.trim();
    const tags     = newTags.value.trim();

    if (!title || !query || !category) {
        showToast('Başlık, sorgu ve kategori zorunludur', 'error');
        return;
    }

    try {
        await AddDork(title, query, category, tags);
        showToast(`"${title}" dork kütüphaneye eklendi`, 'success');
        clearAddForm();
        await loadDorks();
        await refreshStats();
        switchTab('library');
    } catch (e) {
        showToast('Dork eklenirken hata oluştu', 'error');
    }
}

function clearAddForm() {
    newTitle.value    = '';
    newQuery.value    = '';
    newCategory.value = '';
    newTags.value     = '';
    updatePreview();
}

function switchTab(tab) {
    activeTab = tab;
    document.querySelectorAll('.tab-panel').forEach(p => p.classList.remove('active'));
    document.querySelectorAll('.nav-btn').forEach(b => b.classList.remove('active'));

    document.getElementById(`tab-${tab}`).classList.add('active');
    document.querySelector(`.nav-btn[data-tab="${tab}"]`).classList.add('active');

    if (tab === 'history') loadHistory();
}

let toastTimer = null;
function showToast(message, type = 'info') {
    toast.textContent = message;
    toast.className = `toast toast-${type}`;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
        toast.classList.add('hidden');
    }, 3000);
}

function escapeHtml(str) {
    if (!str) return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

function setupEventListeners() {
    document.querySelectorAll('.nav-btn').forEach(btn => {
        btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });

    searchInput.addEventListener('input', () => {
        searchQuery = searchInput.value;
        applyFilters();
    });

    categoryFilter.addEventListener('change', () => {
        activeCategory = categoryFilter.value;
        document.querySelectorAll('.cat-tab').forEach(b => {
            b.classList.toggle('active', b.dataset.cat === activeCategory);
        });
        applyFilters();
    });

    newQuery.addEventListener('input', updatePreview);
    targetInput.addEventListener('input', updatePreview);
    addDorkBtn.addEventListener('click', submitAddDork);
    clearFormBtn.addEventListener('click', clearAddForm);

    clearHistoryBtn.addEventListener('click', async () => {
        if (!confirm('Tüm arama geçmişi silinecek. Emin misiniz?')) return;
        await ClearHistory();
        showToast('Geçmiş temizlendi', 'info');
        await loadHistory();
        await refreshStats();
    });

    [newTitle, newQuery, newCategory, newTags].forEach(input => {
        input.addEventListener('keydown', e => {
            if (e.key === 'Enter') submitAddDork();
        });
    });
}

init();

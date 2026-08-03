// ============================================================
// NasrinShop — админка
// Авторизация: HTTP Basic Auth. Логин/пароль вводятся один раз,
// дальше хранятся в sessionStorage (только для текущей вкладки,
// пропадают при закрытии) и подставляются в каждый запрос.
// ============================================================

const AUTH_KEY = "nasrinshop_admin_auth";

const loginScreen = document.getElementById("login-screen");
const adminPanel = document.getElementById("admin-panel");
const loginForm = document.getElementById("login-form");
const loginError = document.getElementById("login-error");
const logoutBtn = document.getElementById("logout-btn");

const categoriesListEl = document.getElementById("categories-list");
const categorySelectEl = document.getElementById("category-select");
const materialsListEl = document.getElementById("materials-list");
const addCategoryForm = document.getElementById("add-category-form");
const addMaterialForm = document.getElementById("add-material-form");
const messageEl = document.getElementById("admin-message");

function getAuthHeader() {
  const token = sessionStorage.getItem(AUTH_KEY);
  return token ? { Authorization: `Basic ${token}` } : {};
}

function showMessage(text, isError = false) {
  messageEl.textContent = text;
  messageEl.classList.toggle("is-error", isError);
}

// ---------- вход ----------
async function tryLogin(user, pass) {
  const token = btoa(`${user}:${pass}`);
  const res = await fetch(`${API_BASE}/admin/ping`, {
    headers: { Authorization: `Basic ${token}` },
  });
  if (!res.ok) return false;
  sessionStorage.setItem(AUTH_KEY, token);
  return true;
}

loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  loginError.hidden = true;
  const user = document.getElementById("login-user").value.trim();
  const pass = document.getElementById("login-pass").value;

  try {
    const ok = await tryLogin(user, pass);
    if (!ok) {
      loginError.hidden = false;
      return;
    }
    enterPanel();
  } catch (err) {
    console.error(err);
    loginError.textContent = "Сервер недоступен. Проверьте, что бэкенд запущен.";
    loginError.hidden = false;
  }
});

logoutBtn.addEventListener("click", () => {
  sessionStorage.removeItem(AUTH_KEY);
  loginScreen.hidden = false;
  adminPanel.hidden = true;
});

function enterPanel() {
  loginScreen.hidden = true;
  adminPanel.hidden = false;
  loadCategories();
}

// Если логин уже был сделан ранее в этой вкладке — не спрашиваем повторно.
(async function restoreSession() {
  if (!sessionStorage.getItem(AUTH_KEY)) return;
  const res = await fetch(`${API_BASE}/admin/ping`, { headers: getAuthHeader() });
  if (res.ok) enterPanel();
  else sessionStorage.removeItem(AUTH_KEY);
})();

// ---------- хелперы полей (бэкенд может отдавать разный регистр) ----------
function pick(obj, keys, fallback = "") {
  for (const key of keys) {
    if (obj && obj[key] !== undefined && obj[key] !== null) return obj[key];
  }
  return fallback;
}
const getCategoryId = (c) => pick(c, ["id", "ID"]);
const getCategoryName = (c) => pick(c, ["name", "Name"]);
const getMaterialId = (m) => pick(m, ["id", "ID"]);
const getMaterialImage = (m) => pick(m, ["image_url", "ImgURL", "imageUrl", "ImageURL"]);
const getMaterialPrice = (m) => pick(m, ["price", "Price"]);

// ---------- категории ----------
let currentCategories = [];

async function loadCategories() {
  categoriesListEl.innerHTML = "<li>Загрузка…</li>";
  try {
    const res = await fetch(`${API_BASE}/categories`);
    currentCategories = await res.json();
    renderCategories();
    renderCategorySelect();
    if (currentCategories.length > 0) {
      loadMaterials(getCategoryId(currentCategories[0]));
    } else {
      materialsListEl.innerHTML = '<li class="empty-note">Сначала добавь категорию.</li>';
    }
  } catch (err) {
    console.error(err);
    categoriesListEl.innerHTML = "";
    showMessage("Не удалось загрузить категории. Проверьте бэкенд.", true);
  }
}

function renderCategories() {
  categoriesListEl.innerHTML = "";
  if (currentCategories.length === 0) {
    categoriesListEl.innerHTML = '<li class="empty-note">Категорий пока нет.</li>';
    return;
  }
  currentCategories.forEach((cat) => {
    const li = document.createElement("li");
    const name = getCategoryName(cat);
    li.innerHTML = `<span>${escapeHtml(name)}</span>`;
    const delBtn = document.createElement("button");
    delBtn.className = "delete-btn";
    delBtn.type = "button";
    delBtn.textContent = "Удалить";
    delBtn.addEventListener("click", () => deleteCategory(name));
    li.appendChild(delBtn);
    categoriesListEl.appendChild(li);
  });
}

function renderCategorySelect() {
  categorySelectEl.innerHTML = "";
  currentCategories.forEach((cat) => {
    const opt = document.createElement("option");
    opt.value = getCategoryId(cat);
    opt.textContent = getCategoryName(cat);
    categorySelectEl.appendChild(opt);
  });
}

categorySelectEl.addEventListener("change", () => {
  loadMaterials(categorySelectEl.value);
});

addCategoryForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const input = document.getElementById("new-category-name");
  const name = input.value.trim();
  if (!name) return;

  try {
    const res = await fetch(`${API_BASE}/categories/${encodeURIComponent(name)}`, {
      method: "POST",
      headers: getAuthHeader(),
    });
    if (!res.ok) throw new Error(`Статус ${res.status}`);
    input.value = "";
    showMessage(`Категория «${name}» добавлена.`);
    loadCategories();
  } catch (err) {
    console.error(err);
    showMessage("Не удалось добавить категорию.", true);
  }
});

async function deleteCategory(name) {
  if (!confirm(`Удалить категорию «${name}» вместе со всеми материалами?`)) return;
  try {
    const res = await fetch(`${API_BASE}/categories/${encodeURIComponent(name)}`, {
      method: "DELETE",
      headers: getAuthHeader(),
    });
    if (!res.ok) throw new Error(`Статус ${res.status}`);
    showMessage(`Категория «${name}» удалена.`);
    loadCategories();
  } catch (err) {
    console.error(err);
    showMessage("Не удалось удалить категорию.", true);
  }
}

// ---------- материалы ----------
async function loadMaterials(categoryId) {
  if (!categoryId) return;
  materialsListEl.innerHTML = "<li>Загрузка…</li>";
  try {
    const res = await fetch(`${API_BASE}/materials/${categoryId}`);
    const materials = await res.json();
    renderMaterials(materials);
  } catch (err) {
    console.error(err);
    materialsListEl.innerHTML = '<li class="empty-note">Не удалось загрузить материалы.</li>';
  }
}

function renderMaterials(materials) {
  materialsListEl.innerHTML = "";
  if (!materials || materials.length === 0) {
    materialsListEl.innerHTML = '<li class="empty-note">В этой категории пока нет материалов.</li>';
    return;
  }
  materials.forEach((mat) => {
    const li = document.createElement("li");
    const price = getMaterialPrice(mat);
    li.innerHTML = `
      <div class="material-row-info">
        <img class="material-thumb" src="${escapeHtml(getMaterialImage(mat))}" alt="">
        <span>${price ? price + " с." : "без цены"}</span>
      </div>
    `;
    const delBtn = document.createElement("button");
    delBtn.className = "delete-btn";
    delBtn.type = "button";
    delBtn.textContent = "Удалить";
    delBtn.addEventListener("click", () => deleteMaterial(getMaterialId(mat)));
    li.appendChild(delBtn);
    materialsListEl.appendChild(li);
  });
}

addMaterialForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const categoryId = categorySelectEl.value;
  if (!categoryId) {
    showMessage("Сначала выбери или добавь категорию.", true);
    return;
  }
  const imageInput = document.getElementById("new-material-image");
  const priceInput = document.getElementById("new-material-price");

  try {
    const res = await fetch(`${API_BASE}/materials`, {
      method: "POST",
      headers: { "Content-Type": "application/json", ...getAuthHeader() },
      body: JSON.stringify({
        category_id: Number(categoryId),
        price: Number(priceInput.value),
        image_url: imageInput.value.trim(),
      }),
    });
    if (!res.ok) throw new Error(`Статус ${res.status}`);
    imageInput.value = "";
    priceInput.value = "";
    showMessage("Материал добавлен.");
    loadMaterials(categoryId);
  } catch (err) {
    console.error(err);
    showMessage("Не удалось добавить материал.", true);
  }
});

async function deleteMaterial(id) {
  if (!confirm("Удалить этот материал?")) return;
  try {
    const res = await fetch(`${API_BASE}/materials/${id}`, {
      method: "DELETE",
      headers: getAuthHeader(),
    });
    if (!res.ok) throw new Error(`Статус ${res.status}`);
    showMessage("Материал удалён.");
    loadMaterials(categorySelectEl.value);
  } catch (err) {
    console.error(err);
    showMessage("Не удалось удалить материал.", true);
  }
}

// ---------- утилиты ----------
function escapeHtml(str) {
  const div = document.createElement("div");
  div.textContent = str ?? "";
  return div.innerHTML;
}
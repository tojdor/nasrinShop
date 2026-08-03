// ============================================================
// NasrinShop — эндпоинты (адрес бэкенда см. в js/config.js)
// ============================================================
const ENDPOINTS = {
  categories: () => `${API_BASE}/categories`,
  materialsByCategory: (id) => `${API_BASE}/materials/${id}`,
};

// ============================================================
// Тексты интерфейса (RU / TJ). Названия категорий и материалов
// приходят с бэкенда как есть и не переводятся автоматически —
// для перевода карточек товара нужно хранить локали на бэкенде.
// ============================================================
const I18N = {
  ru: {
    tagline: "Ткани и материалы для пошива национальной одежды",
    hero_eyebrow: "NasrinShop",
    hero_title: "Шёлк, атлас и всё для вашего платья",
    hero_body: "Собираем в одном месте ткани, нити и фурнитуру, из которых рождается национальная и современная одежда.",
    error_title: "Не удалось загрузить каталог",
    error_body: "Проверьте соединение с сервером и повторите попытку.",
    retry: "Повторить",
    empty_title: "Коллекции скоро появятся",
    empty_body: "Мы наполняем каталог — загляните немного позже.",
    empty_category: "В этой категории пока нет материалов",
    footer_note: "NasrinShop — материалы для тех, кто шьёт с душой",
  },
  tj: {
    tagline: "Матоъ ва масолеҳ барои дӯхтани либоси миллӣ",
    hero_eyebrow: "NasrinShop",
    hero_title: "Абрешим, атлас ва ҳама чиз барои либоси шумо",
    hero_body: "Мо матоъ, ресмон ва масолеҳи ёрирасонро дар як ҷо ҷамъ мекунем, ки аз онҳо либоси миллӣ ва муосир дӯхта мешавад.",
    error_title: "Каталог бор нашуд",
    error_body: "Пайвастшавӣ ба сервисро тафтиш карда, амалро такрор кунед.",
    retry: "Такрор кардан",
    empty_title: "Коллексияҳо ба зудӣ пайдо мешаванд",
    empty_body: "Мо каталогро пур карда истодаем — каме дертар нигоҳ кунед.",
    empty_category: "Дар ин категория ҳанӯз масолеҳ нест",
    footer_note: "NasrinShop — масолеҳ барои онҳое, ки бо дил медӯзанд",
  },
};

let currentLang = "ru";

function applyI18n(lang){
  currentLang = lang;
  document.documentElement.lang = lang === "tj" ? "tg" : "ru";
  document.documentElement.dataset.lang = lang;
  document.querySelectorAll("[data-i18n]").forEach((el) => {
    const key = el.dataset.i18n;
    if (I18N[lang][key]) el.textContent = I18N[lang][key];
  });
  document.querySelectorAll("[data-lang-btn]").forEach((btn) => {
    btn.classList.toggle("is-active", btn.dataset.langBtn === lang);
  });
}

document.querySelectorAll("[data-lang-btn]").forEach((btn) => {
  btn.addEventListener("click", () => applyI18n(btn.dataset.langBtn));
});

// ============================================================
// Хелперы: бэкенд может отдавать поля в разном регистре
// (id/ID, name/Name и т.д.) — читаем гибко.
// ============================================================
function pick(obj, keys, fallback = "") {
  for (const key of keys) {
    if (obj && obj[key] !== undefined && obj[key] !== null) return obj[key];
  }
  return fallback;
}

function getCategoryId(cat) {
  return pick(cat, ["id", "ID", "Id", "categoryId", "CategoryID"]);
}
function getCategoryName(cat) {
  return pick(cat, ["name", "Name", "title", "Title"], "Без названия");
}
function getMaterialImage(mat) {
  return pick(mat, ["image_url", "imageUrl", "ImageURL", "Image", "image", "url", "URL"]);
}
function getMaterialName(mat) {
  return pick(mat, ["name", "Name", "title", "Title"], "");
}

// ============================================================
// Рендеринг
// ============================================================
const catalogEl = document.getElementById("catalog");
const errorPanel = document.getElementById("error-panel");
const emptyPanel = document.getElementById("empty-panel");
const retryBtn = document.getElementById("retry-btn");

const tplSkeleton = document.getElementById("tpl-skeleton-section");
const tplSection = document.getElementById("tpl-category-section");
const tplCard = document.getElementById("tpl-material-card");

function showSkeletons(count = 3) {
  catalogEl.innerHTML = "";
  for (let i = 0; i < count; i++) {
    catalogEl.appendChild(tplSkeleton.content.cloneNode(true));
  }
}

function clearState() {
  errorPanel.hidden = true;
  emptyPanel.hidden = true;
}

function renderCategorySection(category, materials) {
  const node = tplSection.content.cloneNode(true);
  const titleEl = node.querySelector(".category-title");
  const rowEl = node.querySelector(".material-row");

  titleEl.textContent = getCategoryName(category);

  if (!materials || materials.length === 0) {
    const empty = document.createElement("p");
    empty.className = "category-empty";
    empty.textContent = I18N[currentLang].empty_category;
    rowEl.replaceWith(empty);
    return node;
  }

  if (materials.length <= 4) rowEl.classList.add("is-compact");

  materials.forEach((mat) => {
    const cardNode = tplCard.content.cloneNode(true);
    const img = cardNode.querySelector(".material-img");
    const caption = cardNode.querySelector(".material-name");
    const name = getMaterialName(mat);

    img.src = getMaterialImage(mat);
    img.alt = name || getCategoryName(category);

    if (name) {
      caption.textContent = name;
    } else {
      caption.remove();
    }
    rowEl.appendChild(cardNode);
  });

  return node;
}

// ============================================================
// Загрузка данных
// ============================================================
async function fetchJson(url) {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`Запрос ${url} завершился со статусом ${res.status}`);
  return res.json();
}

async function loadCatalog() {
  clearState();
  showSkeletons();

  try {
    const rawCategories = await fetchJson(ENDPOINTS.categories());
    const categories = Array.isArray(rawCategories) ? rawCategories : [];

    if (categories.length === 0) {
      catalogEl.innerHTML = "";
      emptyPanel.hidden = false;
      return;
    }

    // Материалы по каждой категории запрашиваем параллельно.
    const materialsPerCategory = await Promise.all(
      categories.map((cat) =>
        fetchJson(ENDPOINTS.materialsByCategory(getCategoryId(cat))).catch(() => [])
      )
    );

    catalogEl.innerHTML = "";
    categories.forEach((cat, i) => {
      catalogEl.appendChild(renderCategorySection(cat, materialsPerCategory[i]));
    });
  } catch (err) {
    console.error("Не удалось загрузить каталог NasrinShop:", err);
    catalogEl.innerHTML = "";
    errorPanel.hidden = false;
  }
}

retryBtn.addEventListener("click", loadCatalog);

// ============================================================
// Старт
// ============================================================
applyI18n("ru");
loadCatalog();
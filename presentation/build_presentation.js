const fs = require("fs");
const path = require("path");
const pptxgen = require("pptxgenjs");

const root = path.resolve(__dirname, "..");
const outPath = path.join(__dirname, "Stok_Takip_Proje_Sunumu.pptx");
const assets = path.join(__dirname, "assets");

const pptx = new pptxgen();
pptx.layout = "LAYOUT_WIDE";
pptx.author = "Orkun";
pptx.company = "Go Stok Takip Projesi";
pptx.subject = "Go, PostgreSQL, stok takip ve siparis yonetimi";
pptx.title = "Go Stok Takip Proje Sunumu";
pptx.lang = "tr-TR";
pptx.theme = {
  headFontFace: "Aptos Display",
  bodyFontFace: "Aptos",
  lang: "tr-TR",
};
pptx.defineLayout({ name: "LAYOUT_WIDE", width: 13.333, height: 7.5 });

const C = {
  bg: "101114",
  panel: "171922",
  panel2: "202431",
  text: "F3F6FB",
  muted: "8DA3AA",
  cyan: "00DBE9",
  purple: "9D05FF",
  green: "6EE7B7",
  yellow: "FBBF24",
  red: "F87171",
  border: "2B4650",
};

function imageSize(file) {
  const b = fs.readFileSync(file);
  if (b[0] === 0x89 && b.toString("ascii", 1, 4) === "PNG") {
    return { width: b.readUInt32BE(16), height: b.readUInt32BE(20) };
  }
  return { width: 1440, height: 900 };
}

function addFitImage(slide, file, box) {
  const s = imageSize(file);
  const ratio = s.width / s.height;
  const boxRatio = box.w / box.h;
  let w = box.w;
  let h = box.h;
  if (ratio > boxRatio) {
    h = box.w / ratio;
  } else {
    w = box.h * ratio;
  }
  slide.addImage({
    path: file,
    x: box.x + (box.w - w) / 2,
    y: box.y + (box.h - h) / 2,
    w,
    h,
  });
}

function addTitle(slide, title, subtitle) {
  slide.addShape(pptx.ShapeType.rect, { x: 0, y: 0, w: 13.333, h: 7.5, fill: { color: C.bg }, line: { color: C.bg } });
  slide.addShape(pptx.ShapeType.rect, { x: 0, y: 0, w: 13.333, h: 0.12, fill: { color: C.cyan }, line: { color: C.cyan } });
  slide.addText(title, {
    x: 0.45, y: 0.28, w: 8.8, h: 0.45,
    color: C.text, fontFace: "Aptos Display", fontSize: 22, bold: true,
    margin: 0,
  });
  if (subtitle) {
    slide.addText(subtitle, {
      x: 0.47, y: 0.78, w: 9.5, h: 0.28,
      color: C.muted, fontSize: 9.5, margin: 0,
    });
  }
}

function addPill(slide, text, x, y, w, color = C.cyan) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x, y, w, h: 0.34,
    rectRadius: 0.08,
    fill: { color, transparency: 84 },
    line: { color, transparency: 10 },
  });
  slide.addText(text, { x: x + 0.12, y: y + 0.075, w: w - 0.24, h: 0.16, color, fontSize: 8, bold: true, margin: 0 });
}

function addBullets(slide, items, x, y, w, h, color = C.text) {
  slide.addText(items.map((t) => ({ text: t, options: { bullet: { type: "ul" }, breakLine: true } })), {
    x, y, w, h,
    color,
    fontSize: 10.5,
    breakLine: false,
    fit: "shrink",
    margin: 0.03,
    paraSpaceAfterPt: 7,
  });
}

function addCard(slide, x, y, w, h, title, body, accent = C.cyan) {
  slide.addShape(pptx.ShapeType.roundRect, {
    x, y, w, h,
    rectRadius: 0.08,
    fill: { color: C.panel },
    line: { color: C.border, transparency: 10 },
  });
  slide.addShape(pptx.ShapeType.rect, { x, y, w: 0.08, h, fill: { color: accent }, line: { color: accent } });
  slide.addText(title, { x: x + 0.22, y: y + 0.16, w: w - 0.35, h: 0.22, color: C.text, fontSize: 10.5, bold: true, margin: 0 });
  slide.addText(body, { x: x + 0.22, y: y + 0.48, w: w - 0.35, h: h - 0.58, color: C.muted, fontSize: 8.3, fit: "shrink", margin: 0 });
}

function addSlideNumber(slide, n) {
  slide.addText(String(n).padStart(2, "0"), {
    x: 12.55, y: 7.08, w: 0.42, h: 0.16,
    color: C.muted, fontSize: 8, margin: 0,
  });
}

// Slide 1
{
  const slide = pptx.addSlide();
  addTitle(slide, "Go Stok Takip Projesi", "Admin panelinden urun, kategori, stok ve fiyat yonetimi");
  addPill(slide, "Go net/http", 0.48, 1.22, 1.16);
  addPill(slide, "PostgreSQL", 1.78, 1.22, 1.2, C.green);
  addPill(slide, "Admin Panel", 3.1, 1.22, 1.24, C.purple);
  addPill(slide, "Sepet + Siparis", 4.48, 1.22, 1.45, C.yellow);
  addFitImage(slide, path.join(assets, "admin-products-table.png"), { x: 0.48, y: 1.75, w: 7.9, h: 4.85 });
  slide.addShape(pptx.ShapeType.roundRect, { x: 8.65, y: 1.28, w: 4.15, h: 5.35, rectRadius: 0.08, fill: { color: C.panel }, line: { color: C.border } });
  slide.addText("Admin panelinde eklenen urunler", { x: 8.95, y: 1.58, w: 3.5, h: 0.28, color: C.text, fontSize: 15, bold: true, margin: 0 });
  addBullets(slide, [
    "Telefon: iPhone 15 Pro",
    "Masaustu bilgisayar: Gaming Masaustu Bilgisayar",
    "Laptop: Lenovo IdeaPad Laptop",
    "Stok adedi, fiyat, kategori ve urun gorseli tek panelden yonetiliyor.",
    "Dusuk stok ve tukenme durumlari admin panelinde anlik gorunuyor.",
  ], 8.98, 2.02, 3.45, 2.25);
  addCard(slide, 8.96, 4.66, 3.45, 1.18, "Proje Amaci", "Kucuk bir teknoloji magazasinin urun stoklarini, musteri kayitlarini ve siparis surecini web uzerinden takip etmek.", C.cyan);
  addSlideNumber(slide, 1);
}

// Slide 2
{
  const slide = pptx.addSlide();
  addTitle(slide, "Veritabani Tasarimi ve Is Akisi", "Musteri, adres, sepet, stok ve siparis sureci PostgreSQL tablolarinda tutulur");
  addFitImage(slide, path.join(assets, "admin-customers.png"), { x: 0.48, y: 1.22, w: 5.88, h: 2.78 });
  addFitImage(slide, path.join(assets, "admin-orders.png"), { x: 0.48, y: 4.16, w: 5.88, h: 2.55 });
  addCard(slide, 6.66, 1.18, 2.9, 1.18, "products", "Urun adi, fiyat, stok, kategori, resim ve tanim alanlarini tutar.", C.cyan);
  addCard(slide, 9.72, 1.18, 2.9, 1.18, "categories", "Urunlerin Telefon / Bilgisayar gibi gruplara ayrilmasini saglar.", C.green);
  addCard(slide, 6.66, 2.62, 2.9, 1.18, "musteriler", "Ad, soyad, e-posta, telefon ve sifrelenmis parola bilgisini tutar.", C.purple);
  addCard(slide, 9.72, 2.62, 2.9, 1.18, "adresler", "Musteriye bagli il, ilce, mahalle, acik adres ve varsayilan adres bilgisi.", C.yellow);
  addCard(slide, 6.66, 4.06, 2.9, 1.18, "siparisler", "Musteri, adres, toplam tutar, durum ve siparis tarihini saklar.", C.cyan);
  addCard(slide, 9.72, 4.06, 2.9, 1.18, "siparis_urunleri", "Siparis icindeki urunleri, adetleri ve birim fiyatlari saklar.", C.red);
  slide.addShape(pptx.ShapeType.roundRect, { x: 6.66, y: 5.58, w: 5.96, h: 0.86, rectRadius: 0.08, fill: { color: C.panel2 }, line: { color: C.border } });
  slide.addText("Checkout sirasinda transaction acilir, stok satirlari kilitlenir, stok dusulur ve siparis kaydi olusturulur.", {
    x: 6.95, y: 5.82, w: 5.4, h: 0.26, color: C.text, fontSize: 10.5, bold: true, margin: 0,
  });
  addSlideNumber(slide, 2);
}

pptx.writeFile({ fileName: outPath });

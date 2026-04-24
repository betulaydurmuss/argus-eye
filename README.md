# Argus Eye — Google Dork Aracı

Argus Eye, Google dork sorgularını otomatik olarak oluşturup tarayıcıda çalıştıran masaüstü bir güvenlik aracıdır. Hedef domain veya anahtar kelime girilerek hazır dork kütüphanesinden sorgular seçilir ve Google'da arama yapılır. Kendi dork sorgularını da ekleyip yönetebilirsin.

---

## Özellikler

- 50+ hazır dork sorgusu, 12 kategoride düzenlenmiş
- Hedef domain ile birlikte `site:domain sorgu` formatında otomatik arama
- Kategori ve metin bazlı filtreleme
- Kendi dork sorgularını ekleme ve silme
- Arama geçmişi (son 100 sorgu, kalıcı olarak kaydedilir)
- Geçmişteki sorguları tek tıkla tekrar çalıştırma
- Sorguları panoya kopyalama

---

## Kategoriler

| Kategori | İçerik |
|---|---|
| Yapılandırma & Env | .env, Docker, Kubernetes, Nginx, Git, Firebase, AWS, Azure |
| Veritabanı & Yedek | SQL dump, SQLite, PostgreSQL, MongoDB, Redis, Oracle |
| Admin Panelleri | cPanel, Django, Joomla, SAP, Kibana, Grafana, Outlook |
| Hata & Debug | PHP fatal, Java stack trace, Python traceback, Spring Actuator |
| Hassas Veriler & PII | SSH anahtarı, S3 bucket, Zoom kaydı, kamera akışı, tıbbi kayıt |
| Dizin Listeleme | Index of, /backup, /admin, /uploads, /config, /.git |
| Ağ Cihazları | Cisco WebVPN, FortiGate, Citrix, MikroTik, F5 BIG-IP |
| Doküman Sızıntıları | Confidential PDF, Excel maaş listeleri, CV sızıntıları |
| Web Servis & API | Swagger UI, GraphQL, Postman Collections, Apollo Sandbox |
| Güvenlik Logları | Symantec logs, Firewall kayıtları, SSH login geçmişi |
| Deployment & CI/CD | Terraform State, CircleCI, GitHub Actions, Netlify logs |
| IoT & Endüstriyel | IP Kameralar, Yazıcı panelleri, PLC Web Server, SCADA |

---

## Kurulum

### Gereksinimler

- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Geliştirme modunda çalıştırma

```bash
npm install
wails dev
```

### Derleme (üretim)

```bash
wails build
```

Derlenen uygulama `build/bin/` klasöründe oluşur.

---

## Kullanım

1. Sol paneldeki **Hedef Domain** alanına hedef adresi gir (örn. `example.com`)
2. **Dork Kütüphanesi** sekmesinden bir sorgu seç
3. **Sorgula** butonuna tıkla — Google varsayılan tarayıcında açılır
4. Kendi dorkunu eklemek için **Dork Ekle** sekmesini kullan
5. Geçmiş sorgulara **Geçmiş** sekmesinden ulaşabilirsin

Domain alanı boş bırakılırsa sorgu `site:` öneki olmadan çalışır.

---

## Proje Yapısı

```
argus-eye/
├── app.go              # Backend: dork yönetimi, arama, geçmiş
├── main.go             # Wails uygulama başlangıcı
├── frontend/
│   ├── index.html      # Uygulama arayüzü
│   ├── main.js         # UI mantığı
│   ├── style.css       # Tasarım
│   └── wailsjs/        # Wails tarafından otomatik üretilen JS bağlamaları
└── build/              # Derlenmiş uygulama çıktısı
```

---

## Teknolojiler

- **Go** — Backend mantığı
- **Wails v2** — Go ile masaüstü uygulama çatısı
- **Vanilla JS** — Frontend, harici kütüphane kullanılmadan yazıldı


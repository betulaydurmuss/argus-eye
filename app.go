package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type DorkEntry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Query    string `json:"query"`
	Category string `json:"category"`
	Tags     string `json:"tags"`
}

type HistoryEntry struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	DorkTitle string `json:"dorkTitle"`
	Query     string `json:"query"`
	FullQuery string `json:"fullQuery"`
	Timestamp string `json:"timestamp"`
}

type SearchResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	FullURL string `json:"fullUrl"`
}

type App struct {
	ctx     context.Context
	dorks   []DorkEntry
	history []HistoryEntry
}

func NewApp() *App {
	return &App{
	dorks:   getDefaultDorks(),
	history: []HistoryEntry{},
}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.loadHistory()
}

func (a *App) GetDorks() []DorkEntry {
	return a.dorks
}

func (a *App) GetCategories() []string {
	seen := map[string]bool{}
	var cats []string
	for _, d := range a.dorks {
	if !seen[d.Category] {
	seen[d.Category] = true
	cats = append(cats, d.Category)
}
}
return cats
}

func (a *App) AddDork(title, query, category, tags string) DorkEntry {
newDork := DorkEntry{
	ID:       fmt.Sprintf("custom_%d", time.Now().Unix()),
	Title:    title,
	Query:    query,
	Category: category,
	Tags:     tags,
}
a.dorks = append(a.dorks, newDork)
return newDork
}

func (a *App) DeleteDork(id string) bool {
for i, d := range a.dorks {
	if d.ID == id {
	a.dorks = append(a.dorks[:i], a.dorks[i+1:]...)
	return true
}
}
return false
}

func (a *App) RunGoogleDork(domain, dorkTitle, dorkQuery string) SearchResult {
var fullQuery string
if domain != "" {
	fullQuery = fmt.Sprintf("site:%s %s", domain, dorkQuery)
} else {
	fullQuery = dorkQuery
}

searchURL := "https://www.google.com/search?q=" + url.QueryEscape(fullQuery)

var err error
switch runtime.GOOS {
case "windows":
	err = exec.Command("rundll32", "url.dll,FileProtocolHandler", searchURL).Start()
case "darwin":
	err = exec.Command("open", searchURL).Start()
default: 
	err = exec.Command("xdg-open", searchURL).Start()
}

if err != nil {
	return SearchResult{
	Success: false,
	Message: "Tarayıcı açılamadı: " + err.Error(),
	FullURL: searchURL,
}
}

record := HistoryEntry{
	ID:        fmt.Sprintf("h_%d", time.Now().Unix()),
	Domain:    domain,
	DorkTitle: dorkTitle,
	Query:     dorkQuery,
	FullQuery: fullQuery,
	Timestamp: time.Now().Format("02.01.2006 15:04:05"),
}
a.history = append([]HistoryEntry{record}, a.history...)
if len(a.history) > 100 {
a.history = a.history[:100]
}
a.saveHistory()

return SearchResult{Success: true, Message: "Tarayıcıda açıldı", FullURL: searchURL}
}

func (a *App) GetHistory() []HistoryEntry {
	return a.history
}

func (a *App) ClearHistory() bool {
	a.history = []HistoryEntry{}
	a.saveHistory()
	return true
}

func (a *App) GetStats() map[string]int {
	cats := map[string]bool{}
	for _, d := range a.dorks {
	cats[d.Category] = true
}
return map[string]int{
	"totalDorks":      len(a.dorks),
	"totalHistory":    len(a.history),
	"totalCategories": len(cats),
}
}

func historyFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".argus-eye-history.json")
}

func (a *App) saveHistory() {
	data, err := json.Marshal(a.history)
	if err != nil {
	return
}
_ = os.WriteFile(historyFilePath(), data, 0644)
}

func (a *App) loadHistory() {
	data, err := os.ReadFile(historyFilePath())
	if err != nil {
	return
}
_ = json.Unmarshal(data, &a.history)
}

func getDefaultDorks() []DorkEntry {
	return []DorkEntry{
		{ID: "e001", Title: "Ortam Değişkeni Sızıntısı", Query: "filetype:env \"DB_PASSWORD\"", Category: "Yapılandırma & Env", Tags: "env,secrets,database"},
		{ID: "e002", Title: "Docker Compose Dosyası", Query: "filename:docker-compose.yml \"root\"", Category: "Yapılandırma & Env", Tags: "docker,compose,root"},
		{ID: "e003", Title: "Kubernetes Secret", Query: "extension:yaml \"kind: Secret\"", Category: "Yapılandırma & Env", Tags: "kubernetes,k8s,secret"},
		{ID: "e004", Title: "Nginx Yapılandırması", Query: "filetype:conf inurl:nginx.conf", Category: "Yapılandırma & Env", Tags: "nginx,config,server"},
		{ID: "e005", Title: "Git Config Sızıntısı", Query: "inurl:/.git/config", Category: "Yapılandırma & Env", Tags: "git,config,leak"},
		{ID: "e006", Title: "Node.js Bağımlılıkları", Query: "filename:package-lock.json \"dependencies\"", Category: "Yapılandırma & Env", Tags: "nodejs,npm,dependencies"},
		{ID: "e007", Title: "WordPress Yapılandırması", Query: "filetype:php \"wp-config.php\" -github", Category: "Yapılandırma & Env", Tags: "wordpress,config,php"},
		{ID: "e008", Title: "Azure Bağlantı Dizesi", Query: "filetype:config \"connectionString\"", Category: "Yapılandırma & Env", Tags: "azure,connection,string"},
		{ID: "e009", Title: "AWS Kimlik Bilgileri", Query: "filetype:ini \"aws_access_key_id\"", Category: "Yapılandırma & Env", Tags: "aws,credentials,cloud"},
		{ID: "e010", Title: "Firebase Yapılandırması", Query: "filetype:json \"firebaseConfig\"", Category: "Yapılandırma & Env", Tags: "firebase,json,config"},

		{ID: "f001", Title: "SQL Dump Dosyası", Query: "extension:sql \"INSERT INTO\"", Category: "Veritabanı & Yedek", Tags: "sql,dump,database"},
		{ID: "f002", Title: "SQLite Veritabanı", Query: "extension:db OR extension:sqlite", Category: "Veritabanı & Yedek", Tags: "sqlite,db,database"},
		{ID: "f003", Title: "Yedek Arşiv Dosyası", Query: "filetype:zip OR filetype:tar.gz \"backup\"", Category: "Veritabanı & Yedek", Tags: "backup,zip,archive"},
		{ID: "f004", Title: "PostgreSQL Log", Query: "extension:log \"postgresql\"", Category: "Veritabanı & Yedek", Tags: "postgresql,log,database"},
		{ID: "f005", Title: "MongoDB Geçmişi", Query: "filename:.dbshell", Category: "Veritabanı & Yedek", Tags: "mongodb,shell,history"},
		{ID: "f006", Title: "Redis Yapılandırması", Query: "extension:conf \"redis.conf\"", Category: "Veritabanı & Yedek", Tags: "redis,config,cache"},
		{ID: "f007", Title: "Oracle TNS Dosyası", Query: "filename:tnsnames.ora", Category: "Veritabanı & Yedek", Tags: "oracle,tns,database"},
		{ID: "f008", Title: "Veritabanı Migrasyon", Query: "inurl:/migrations/ extension:sql", Category: "Veritabanı & Yedek", Tags: "migration,sql,database"},
		{ID: "f009", Title: "Access Veritabanı", Query: "extension:mdb OR extension:accdb", Category: "Veritabanı & Yedek", Tags: "access,mdb,database"},
		{ID: "f010", Title: "Log Dosyasında Şifre", Query: "filetype:log \"password\" OR \"login\"", Category: "Veritabanı & Yedek", Tags: "log,password,leak"},

		{ID: "g001", Title: "Kontrol Paneli", Query: "intitle:\"Control Panel\" inurl:admin", Category: "Admin Panelleri", Tags: "control,panel,admin"},
		{ID: "g002", Title: "cPanel Arayüzü", Query: "inurl:2083 OR inurl:2082", Category: "Admin Panelleri", Tags: "cpanel,hosting,panel"},
		{ID: "g003", Title: "Laravel Debug Açık", Query: "intext:\"APP_DEBUG=true\"", Category: "Admin Panelleri", Tags: "laravel,debug,php"},
		{ID: "g004", Title: "Django Yönetim Paneli", Query: "inurl:/admin/login/", Category: "Admin Panelleri", Tags: "django,admin,python"},
		{ID: "g005", Title: "Joomla Yönetim Paneli", Query: "inurl:/administrator/index.php", Category: "Admin Panelleri", Tags: "joomla,cms,admin"},
		{ID: "g006", Title: "SAP Giriş Sayfası", Query: "intitle:\"SAP NetWeaver Portal\"", Category: "Admin Panelleri", Tags: "sap,netweaver,enterprise"},
		{ID: "g007", Title: "Kibana Paneli", Query: "inurl:5601/app/kibana", Category: "Admin Panelleri", Tags: "kibana,elasticsearch,dashboard"},
		{ID: "g008", Title: "Grafana Giriş", Query: "intitle:\"Grafana\" inurl:/login", Category: "Admin Panelleri", Tags: "grafana,monitoring,dashboard"},
		{ID: "g009", Title: "Outlook Web Erişimi", Query: "inurl:/owa/auth/login.aspx", Category: "Admin Panelleri", Tags: "outlook,owa,exchange"},
		{ID: "g010", Title: "SonicWall Giriş", Query: "inurl:/cgi-bin/welcome", Category: "Admin Panelleri", Tags: "sonicwall,firewall,vpn"},

		{ID: "h001", Title: "SQL Sözdizimi Hatası", Query: "intext:\"sql syntax near\"", Category: "Hata & Debug", Tags: "sql,error,syntax"},
		{ID: "h002", Title: "PHP Fatal Hatası", Query: "intext:\"Fatal error:\" \"on line\"", Category: "Hata & Debug", Tags: "php,fatal,error"},
		{ID: "h003", Title: "ASP.NET Hata Sayfası", Query: "intext:\"Runtime Error\" \"Stack Trace:\"", Category: "Hata & Debug", Tags: "aspnet,runtime,stacktrace"},
		{ID: "h004", Title: "Java Stack Trace", Query: "intext:\"at java.lang.\" \"Exception in thread\"", Category: "Hata & Debug", Tags: "java,exception,stacktrace"},
		{ID: "h005", Title: "Python Hata İzi", Query: "intext:\"Traceback (most recent call last):\"", Category: "Hata & Debug", Tags: "python,traceback,error"},
		{ID: "h006", Title: "ColdFusion Hatası", Query: "extension:cfm \"Error Occurred While Processing Request\"", Category: "Hata & Debug", Tags: "coldfusion,cfm,error"},
		{ID: "h007", Title: "Rails Debug Sayfası", Query: "intext:\"Extracted source (around line #)\"", Category: "Hata & Debug", Tags: "rails,ruby,debug"},
		{ID: "h008", Title: "XAMPP PHP Bilgisi", Query: "inurl:/xampp/phpinfo.php", Category: "Hata & Debug", Tags: "xampp,phpinfo,php"},
		{ID: "h009", Title: "Symantec Log Dosyası", Query: "filetype:log \"Symantec Endpoint Protection\"", Category: "Hata & Debug", Tags: "symantec,antivirus,log"},
		{ID: "h010", Title: "Spring Boot Actuator", Query: "inurl:/actuator/env", Category: "Hata & Debug", Tags: "spring,actuator,java"},

		{ID: "j001", Title: "Admin Login Sayfası", Query: "inurl:admin login", Category: "Admin & Giriş", Tags: "admin,login"},
		{ID: "j002", Title: "Admin Login (intitle)", Query: "intitle:\"admin login\"", Category: "Admin & Giriş", Tags: "admin,login,title"},
		{ID: "j003", Title: "Login PHP", Query: "inurl:login.php", Category: "Admin & Giriş", Tags: "login,php"},
		{ID: "j004", Title: "WordPress Admin", Query: "inurl:wp-admin", Category: "Admin & Giriş", Tags: "wordpress,admin"},
		{ID: "j005", Title: "cPanel Girişi", Query: "inurl:cpanel", Category: "Admin & Giriş", Tags: "cpanel,hosting"},
		{ID: "j006", Title: "Dashboard Admin", Query: "intitle:\"Dashboard\" inurl:admin", Category: "Admin & Giriş", Tags: "dashboard,admin"},
		{ID: "j007", Title: "Admin Panel URL", Query: "inurl:adminpanel", Category: "Admin & Giriş", Tags: "adminpanel"},
		{ID: "j008", Title: "Backend Login", Query: "inurl:backend login", Category: "Admin & Giriş", Tags: "backend,login"},
		{ID: "j009", Title: "Admin Panel (intitle)", Query: "intitle:\"Admin Panel\"", Category: "Admin & Giriş", Tags: "admin,panel,title"},
		{ID: "j010", Title: "Signin Sayfası", Query: "inurl:signin", Category: "Admin & Giriş", Tags: "signin,login"},

		{ID: "k001", Title: "Backup Dizini", Query: "intitle:\"index of\" /backup", Category: "Dizin Listeleme", Tags: "backup,directory"},
		{ID: "k002", Title: "Admin Dizini", Query: "intitle:\"index of\" /admin", Category: "Dizin Listeleme", Tags: "admin,directory"},
		{ID: "k003", Title: "Password Dizini", Query: "intitle:\"index of\" /password", Category: "Dizin Listeleme", Tags: "password,directory"},
		{ID: "k004", Title: "Git Dizini", Query: "intitle:\"index of\" .git", Category: "Dizin Listeleme", Tags: "git,directory"},
		{ID: "k005", Title: "Veritabanı Dizini", Query: "intitle:\"index of\" \"database\"", Category: "Dizin Listeleme", Tags: "database,directory"},
		{ID: "k006", Title: "Uploads Dizini", Query: "intitle:\"index of\" /uploads", Category: "Dizin Listeleme", Tags: "uploads,directory"},
		{ID: "k007", Title: "Config Dizini", Query: "intitle:\"index of\" config", Category: "Dizin Listeleme", Tags: "config,directory"},
		{ID: "k008", Title: "Env Dizini", Query: "intitle:\"index of\" .env", Category: "Dizin Listeleme", Tags: "env,directory"},

		{ID: "l001", Title: "SQL Şifre Dump", Query: "filetype:sql \"password\"", Category: "Hassas Dosyalar", Tags: "sql,password,dump"},
		{ID: "l002", Title: "Env DB Şifresi", Query: "filetype:env DB_PASSWORD", Category: "Hassas Dosyalar", Tags: "env,db,password"},
		{ID: "l003", Title: "Log Hata Dosyası", Query: "filetype:log \"error\"", Category: "Hassas Dosyalar", Tags: "log,error"},
		{ID: "l004", Title: "TXT Kullanıcı Adı Şifre", Query: "filetype:txt \"username\" \"password\"", Category: "Hassas Dosyalar", Tags: "txt,username,password"},
		{ID: "l005", Title: "Backup BAK Dosyası", Query: "filetype:bak inurl:backup", Category: "Hassas Dosyalar", Tags: "bak,backup"},
		{ID: "l006", Title: "XML Config", Query: "filetype:xml \"config\"", Category: "Hassas Dosyalar", Tags: "xml,config"},
		{ID: "l007", Title: "JSON API Anahtarı", Query: "filetype:json \"api_key\"", Category: "Hassas Dosyalar", Tags: "json,api,key"},
		{ID: "l008", Title: "INI Şifre", Query: "filetype:ini \"password\"", Category: "Hassas Dosyalar", Tags: "ini,password"},

		{ID: "m001", Title: "PHP Info Sayfası", Query: "inurl:phpinfo.php", Category: "Sunucu & Sistem", Tags: "phpinfo,php,server"},
		{ID: "m002", Title: "Apache Durum", Query: "intitle:\"Apache Status\"", Category: "Sunucu & Sistem", Tags: "apache,status,server"},
		{ID: "m003", Title: "Apache Index", Query: "intitle:\"Index of /\" \"server at\"", Category: "Sunucu & Sistem", Tags: "apache,index,server"},
		{ID: "m004", Title: "Test PHP Dosyası", Query: "inurl:test.php", Category: "Sunucu & Sistem", Tags: "test,php"},
		{ID: "m005", Title: "Debug Sayfası", Query: "inurl:debug", Category: "Sunucu & Sistem", Tags: "debug"},
		{ID: "m006", Title: "Env Endpoint", Query: "inurl:env", Category: "Sunucu & Sistem", Tags: "env,endpoint"},

		{ID: "n001", Title: "Kamera ViewerFrame", Query: "inurl:\"ViewerFrame?Mode=\"", Category: "Kamera & IoT", Tags: "camera,iot,viewer"},
		{ID: "n002", Title: "AXIS Kamera", Query: "intitle:\"Live View / - AXIS\"", Category: "Kamera & IoT", Tags: "axis,camera,live"},
		{ID: "n003", Title: "Kamera view.shtml", Query: "inurl:/view.shtml", Category: "Kamera & IoT", Tags: "camera,shtml"},
		{ID: "n004", Title: "WebcamXP", Query: "intitle:\"webcamXP\"", Category: "Kamera & IoT", Tags: "webcam,iot"},
		{ID: "n005", Title: "Video CGI Akışı", Query: "inurl:/video.cgi", Category: "Kamera & IoT", Tags: "video,cgi,camera"},

		{ID: "o001", Title: "TXT Şifre Parametresi", Query: "intext:\"password=\" filetype:txt", Category: "Kimlik Sızıntısı", Tags: "password,txt,leak"},
		{ID: "o002", Title: "DB_PASSWORD Metni", Query: "intext:\"DB_PASSWORD\"", Category: "Kimlik Sızıntısı", Tags: "db,password,leak"},
		{ID: "o003", Title: "API Key Metni", Query: "intext:\"api_key\"", Category: "Kimlik Sızıntısı", Tags: "api,key,leak"},
		{ID: "o004", Title: "Secret Key Metni", Query: "intext:\"secret_key\"", Category: "Kimlik Sızıntısı", Tags: "secret,key,leak"},
		{ID: "o005", Title: "Bearer Token", Query: "intext:\"Authorization: Bearer\"", Category: "Kimlik Sızıntısı", Tags: "bearer,token,auth"},
		{ID: "o006", Title: "AWS Access Key", Query: "intext:\"aws_access_key_id\"", Category: "Kimlik Sızıntısı", Tags: "aws,access,key"},

		{ID: "p001", Title: "Site Admin", Query: "inurl:admin", Category: "Site Bazlı", Tags: "site,admin"},
		{ID: "p002", Title: "Site SQL Dosyası", Query: "filetype:sql", Category: "Site Bazlı", Tags: "site,sql"},
		{ID: "p003", Title: "Site Backup", Query: "inurl:backup", Category: "Site Bazlı", Tags: "site,backup"},
		{ID: "p004", Title: "Site Log Dosyası", Query: "ext:log", Category: "Site Bazlı", Tags: "site,log,ext"},
		{ID: "p005", Title: "Site Dizin Listesi", Query: "\"index of\"", Category: "Site Bazlı", Tags: "site,index,directory"},

		{ID: "i001", Title: "Gizli PDF Belgesi", Query: "filetype:pdf \"confidential\" OR \"not for distribution\"", Category: "Hassas Veriler & PII", Tags: "pdf,confidential,document"},
		{ID: "i002", Title: "Excel'de Şifre", Query: "filetype:xlsx \"password\" OR \"identity\"", Category: "Hassas Veriler & PII", Tags: "excel,password,identity"},
		{ID: "i003", Title: "Çalışan Listesi", Query: "filetype:csv \"employee_name\" \"salary\"", Category: "Hassas Veriler & PII", Tags: "csv,employee,salary,pii"},
		{ID: "i004", Title: "Tıbbi Kayıt", Query: "filetype:pdf \"medical record\" OR \"patient info\"", Category: "Hassas Veriler & PII", Tags: "medical,patient,record,pii"},
		{ID: "i005", Title: "VPN Yapılandırma Dosyası", Query: "extension:ovpn OR extension:pcf", Category: "Hassas Veriler & PII", Tags: "vpn,ovpn,config"},
		{ID: "i006", Title: "SSH Özel Anahtarı", Query: "filename:id_rsa OR filename:id_dsa", Category: "Hassas Veriler & PII", Tags: "ssh,private,key,rsa"},
		{ID: "i007", Title: "Pastebin Kimlik Dökümü", Query: "site:pastebin.com \"admin:\"", Category: "Hassas Veriler & PII", Tags: "pastebin,dump,credentials"},
		{ID: "i008", Title: "S3 Bucket Dizini", Query: "site:s3.amazonaws.com \"index of /\"", Category: "Hassas Veriler & PII", Tags: "s3,aws,bucket,cloud"},
		{ID: "i009", Title: "Zoom Toplantı Kaydı", Query: "inurl:zoom.us/rec/play", Category: "Hassas Veriler & PII", Tags: "zoom,recording,meeting"},
		{ID: "i010", Title: "Ağ Kamerası Yayını", Query: "inurl:\"view.shtml\" intitle:\"network camera\"", Category: "Hassas Veriler & PII", Tags: "camera,cctv,feed,network"},
	}
}

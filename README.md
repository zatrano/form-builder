# go-form

<p align="center">
  <img src="https://raw.githubusercontent.com/gofiber/logo/master/rebranding-logo-with-name-dark.png" alt="Go-Form Logo Placeholder" width="300">
</p>

<p align="center">
  A flexible, secure, and component-friendly form management package for Go, featuring built-in CSRF protection and powerful validation.
  <br>
  <a href="https://pkg.go.dev/github.com/zatrano/go-form"><strong>GoDoc</strong></a>
  ·
  <a href="https://github.com/zatrano/go-form/issues">Report Bug</a>
  ·
  <a href="https://github.com/zatrano/go-form/issues">Request Feature</a>
</p>

---

`go-form`, sunucu taraflı render (SSR) yapan Go web uygulamaları için modern, güvenli ve yeniden kullanılabilir formlar oluşturmayı basitleştirmek üzere tasarlanmıştır. Karmaşık HTML'i, validasyon mantığını ve güvenlik endişelerini soyutlayarak, geliştiricinin sadece iş mantığına odaklanmasını sağlar.

## ✨ Özellikler

*   **Güvenli:** Dahili, depolama (storage) bağımsız CSRF koruması.
*   **Güçlü Validasyon:** Endüstri standardı `go-playground/validator/v10` ile kutudan çıktığı gibi entegrasyon.
*   **Bileşen Dostu:** Yeniden kullanılabilir HTML bileşenleri (input, select, textarea vb.) oluşturmayı destekleyen yardımcı fonksiyonlar.
*   **Hata ve Durum Yönetimi:** Validasyon hatalarını ve kullanıcının eski girdilerini (`old input`) kolayca yönetin ve view'a geri gönderin.
*   **Framework Bağımsız:** Çekirdek mantık herhangi bir framework ile çalışır. [Fiber](https://gofiber.io/) için kullanıma hazır bir adaptör içerir.
*   **Esnek Depolama:** CSRF token'ları için kendi depolama mekanizmanızı (session, cookie, Redis vb.) kolayca bağlayın.

## 🚀 Kurulum

```bash
go get github.com/zatrano/go-form
```

## 📋 Hızlı Başlangıç (Fiber ile)

Bu örnek, `go-form`'un bir Fiber web sunucusunda nasıl kullanılacağını gösterir.

### 1. Form Struct'ını Tanımla

Uygulamanızın `forms` klasöründe, `validate` etiketlerini kullanarak formunuzu tanımlayın.

```go
// app/forms/register_form.go
package forms

type RegisterForm struct {
    Name     string `form:"name" validate:"required,min=3"`
    Email    string `form:"email" validate:"required,email"`
    Password string `form:"password" validate:"required,min=8"`
}
```

### 2. View Helper'larını Ekle

Template motorunuzu başlatırken, `go-form`'un yardımcı fonksiyonlarını ekleyin.

```go
// internal/zatrano/view/engine.go
import (
    formhelper "github.com/zatrano/go-form"
    "github.com/gofiber/template/html"
)

func NewEngine() *html.Engine {
    engine := html.New("./views", ".html")

    // Form helper fonksiyonlarını motorumuza kaydediyoruz.
    for name, fn := range formhelper.GetTemplateFuncs() {
        engine.AddFunc(name, fn)
    }
    // ...
    return engine
}
```

### 3. Handler'da Kullanım

Handler'larınızda GET istekleri için boş bir form oluşturun ve POST isteklerini işleyin.

```go
// app/handlers/auth_handler.go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/session"
    form "github.com/zatrano/go-form"
    "github.com/zatrano/go-form/adapter"
)

// ShowRegisterForm, kayıt formunu gösterir (GET)
func (h *AuthHandler) ShowRegisterForm(c *fiber.Ctx) error {
    store := c.Locals("session_store").(*session.Store)
    
    // Yeni bir form oluştur, bu işlem CSRF token'ını üretip session'a kaydeder.
    formData, _ := adapter.NewFromContext(store, c)

    return c.Render("register", fiber.Map{"Form": formData})
}

// HandleRegisterForm, kayıt formunu işler (POST)
func (h *AuthHandler) HandleRegisterForm(c *fiber.Ctx) error {
    store := c.Locals("session_store").(*session.Store)
    sess, _ := store.Get(c)
    
    var model forms.RegisterForm
    c.BodyParser(&model)
    formValues, _ := c.FormValues()

    // Formu parse et: CSRF'yi doğrular ve validasyonu çalıştırır.
    formData, err := form.Parse(form.Config{
        Storage:   adapter.NewFiberStorage(store, c),
        SessionID: sess.ID(),
        CSRFField: "_csrf",
    }, formValues, &model)

    if err != nil {
        // Hata varsa, hataları ve eski girdileri içeren formu tekrar render et.
        return c.Render("register", fiber.Map{"Form": formData})
    }
    
    // Validasyon başarılı! Kullanıcıyı kaydet...
    // ...
    
    return c.Redirect("/dashboard")
}
```

### 4. View'da Form Oluşturma

HTML template'lerinizde, `Form` nesnesini kullanarak yeniden kullanılabilir bileşenler oluşturun.

```html
<!-- views/register.html -->
<form action="/register" method="POST">
    <!-- CSRF Token'ını ekle -->
    <input type="hidden" name="{{.Form.CSRFField}}" value="{{.Form.CSRFToken}}">

    <!-- Name alanı -->
    <div class="mb-3">
        <label>Full Name</label>
        <input 
            type="text" 
            name="name"
            class="form-control {{if hasError "name" .Form.Errors}}is-invalid{{end}}"
            value="{{old "name" .Form.OldInput}}"
        >
        {{if hasError "name" .Form.Errors}}
            <div class="invalid-feedback">{{getError "name" .Form.Errors}}</div>
        {{end}}
    </div>

    <!-- ... diğer form alanları ... -->

    <button type="submit">Register</button>
</form>
```

## 🛠️ Gelişmiş Kullanım

### Kendi Depolama Adaptörünü Yazma

`go-form`, CSRF token'larını saklamak için bir `Storage` arayüzü kullanır. Fiber dışındaki bir framework veya özel bir session yönetimi için kendi adaptörünüzü kolayca yazabilirsiniz.

```go
type MyCustomStorage struct {
    // ...
}

func (s *MyCustomStorage) Get(sessionID string) string {
    // ... Get token logic ...
}

func (s *MyCustomStorage) Set(sessionID, token string) {
    // ... Set token logic ...
}

func (s *MyCustomStorage) Delete(sessionID string) {
    // ... Delete token logic ...
}
```

## 🤝 Katkıda Bulunma

Katkılarınız projeyi daha da iyi hale getirir! Lütfen bir "issue" açın veya "pull request" gönderin.

1.  Projeyi Fork'layın (`https://github.com/zatrano/go-form/fork`)
2.  Yeni bir Feature Branch'i oluşturun (`git checkout -b feature/AmazingFeature`)
3.  Değişikliklerinizi Commit'leyin (`git commit -m 'Add some AmazingFeature'`)
4.  Branch'i Push'layın (`git push origin feature/AmazingFeature`)
5.  Bir Pull Request açın.

## 📄 Lisans

Bu proje MIT Lisansı altında dağıtılmaktadır. Daha fazla bilgi için `LICENSE` dosyasına bakın.

---
zatrano kısımlarını kendi GitHub kullanıcı adınızla değiştirmeyi unutmayın.

package main

import (
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// WhatsApp OTP sistemi için gerekli ayarlar
	app.OnBootstrap().Bind(&core.Handler[*core.BootstrapEvent]{
		Func: func(e *core.BootstrapEvent) error {
			// WhatsApp Business API ayarlarını yapılandır
			settings := e.App.Settings()
			settings.Meta.WhatsAppAccessToken = "YOUR_WHATSAPP_ACCESS_TOKEN"
			settings.Meta.WhatsAppPhoneNumberID = "YOUR_PHONE_NUMBER_ID"
			
			// Collection OTP ayarlarını yapılandır
			collection, err := e.App.FindCollectionByNameOrId("users")
			if err == nil {
				collection.OTP.Enabled = true
				collection.OTP.DeliveryMethod = "whatsapp" // "email", "whatsapp", "both"
				collection.OTP.Duration = 300 // 5 dakika
				collection.OTP.Length = 6
				
				// WhatsApp mesaj template'ini özelleştir
				collection.OTP.WhatsAppTemplate.Message = `Merhaba! Doğrulama kodunuz: *{{OTP}}*

Bu kodu kimseyle paylaşmayın. Kod 5 dakika geçerlidir.

Teşekkürler,
{{APP_NAME}} ekibi`
				
				e.App.Save(collection)
			}

			return e.Next()
		},
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

/*
WhatsApp Business API Kurulumu:

1. Meta for Developers'da WhatsApp Business API hesabı oluşturun
2. Access Token ve Phone Number ID'nizi alın
3. Yukarıdaki örnekte bu bilgileri güncelleyin

API Kullanımı:

1. OTP Talep Etme:
   POST /api/collections/users/request-otp
   {
     "email": "user@example.com"
   }

2. OTP ile Giriş:
   POST /api/collections/users/auth-with-otp
   {
     "otpId": "otp_id_from_request",
     "password": "123456"
   }

Önemli Notlar:
- Kullanıcı kaydında "phone" alanı bulunmalı
- WhatsApp mesajları sadece kayıtlı numaralara gönderilir
- Rate limiting koruması mevcuttur
- Email ve WhatsApp birlikte kullanılabilir
*/

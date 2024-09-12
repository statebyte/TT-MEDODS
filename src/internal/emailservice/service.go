package emailservice

import (
	"log"
)

func SendWarning(email string) {
	// Моковая отправка email
	log.Println("[MOCK] Warning: IP address changed for user:", email)
}

package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityMiddleware adiciona headers de segurança
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Previne clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Proteção XSS
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Impede que o navegador faça MIME-sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Política de segurança de conteúdo
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; img-src 'self' data:; style-src 'self'; font-src 'self'; connect-src 'self'")
		
		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Strict Transport Security (força HTTPS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Impede que o site seja incorporado em iframes
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		
		c.Next()
	}
}

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
)

var (
	urlStore    = make(map[string]string)
	secretKey   = "shhhhh_this_is_an_dumb_key123456" // 32 bytes for AES-256
	mu          sync.Mutex
	lettersRune = []rune("abcdefghijlmnopqrstuvwxyzABCDEFGHIJLMNOPQRSTUVWXYZ123456789")
)

func decrypt(encryptedUrl string) string {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}

	cipherText, err := hex.DecodeString(encryptedUrl)

	if err != nil {
		log.Fatal(err)
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func encrypt(originalUrl string) string {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		log.Fatal(err)
	}
	plainText := []byte(originalUrl)
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
	return hex.EncodeToString(cipherText)
}

func generateShortId() string {
	b := make([]rune, 6)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			log.Fatal(err)
		}
		b[i] = lettersRune[num.Int64()]
	}
	return string(b)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortId := r.URL.Path[1:]

	// Don't handle empty paths
	if shortId == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	mu.Lock()
	encryptedUrl, ok := urlStore[shortId]
	mu.Unlock()

	if !ok {
		http.Error(w, "This URL doesnt exist", http.StatusNotFound)
		return
	}

	decryptedUrl := decrypt(encryptedUrl)
	http.Redirect(w, r, decryptedUrl, http.StatusFound)
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	originalUrl := r.URL.Query().Get("url")
	if originalUrl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	encryptedUrl := encrypt(originalUrl)
	mu.Lock()
	shortId := generateShortId()
	urlStore[shortId] = encryptedUrl
	mu.Unlock()
	shortUrl := fmt.Sprintf("http://localhost:8080/%s", shortId)
	fmt.Fprintf(w, "The shorted URL of this original URL is: %s", shortUrl)
}

func main() {
	http.HandleFunc("/shorten", shortenUrl)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("Running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

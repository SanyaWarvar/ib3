package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// EncryptStruct encrypts all fields of the given struct.
func EncryptStruct(s interface{}, key []byte) error {
	v := reflect.ValueOf(s).Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		// Сериализуем значение поля в JSON
		data := field.Interface().(string)

		// Шифруем сериализованные данные
		ciphertext, err := Encrypt([]byte(data), key)
		if err != nil {
			return err
		}

		strValue := hex.EncodeToString(ciphertext)
		field.SetString(strValue)
	}
	return nil
}

// DecryptStruct decrypts all fields of the given struct.
func DecryptStruct(s interface{}, key []byte) error {
	v := reflect.ValueOf(s).Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue // Пропускаем поля, которые нельзя установить
		}

		// Декодируем зашифрованные данные
		ciphertext := field.Interface().(string)
		data, _ := hex.DecodeString(ciphertext)
		plaintext, err := Decrypt(data, key)
		if err != nil {
			return err
		}

		// Устанавливаем значение обратно в поле
		field.SetString(string(plaintext))
	}
	return nil
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptData расшифровывает данные с использованием AES и возвращает исходные данные в виде []byte.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

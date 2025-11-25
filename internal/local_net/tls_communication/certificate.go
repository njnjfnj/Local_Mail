package tls_communication

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func GetOrGenerateCertificate(certFile, keyFile string) (tls.Certificate, error) {
	// 1. Попытка загрузить существующие сертификаты
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err == nil {
		return cert, nil // Успешно загрузили старые
	}

	// 2. Если файлов нет или ошибка — генерируем новые
	// Генерируем приватный ключ RSA (2048 бит)
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Настраиваем шаблон сертификата
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Валиден 1 год

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"My P2P Messenger"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		// Важно: разрешаем использование и для сервера, и для клиента
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// Создаем сам сертификат (байты)
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 3. Сохраняем в файлы (PEM кодировка)
	// Сохраняем cert.pem
	certOut, err := os.Create(certFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// Сохраняем key.pem
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	keyOut.Close()

	// 4. Загружаем и возвращаем то, что только что создали
	return tls.LoadX509KeyPair(certFile, keyFile)
}

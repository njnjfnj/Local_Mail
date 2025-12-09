package tls_communication

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"fyne.io/fyne/v2"
)

const (
	PREF_CERT = "tls_cert_pem_v2"
	PREF_KEY  = "tls_key_pem_v2"
)

func GetOrGenerateCertificate(app fyne.App) (tls.Certificate, error) {
	certStr := app.Preferences().StringWithFallback(PREF_CERT, "")
	keyStr := app.Preferences().StringWithFallback(PREF_KEY, "")

	if certStr != "" && keyStr != "" {
		cert, err := tls.X509KeyPair([]byte(certStr), []byte(keyStr))
		if err == nil {
			return cert, nil
		}
	}

	certPEM, keyPEM, err := generatePEMKeys()
	if err != nil {
		return tls.Certificate{}, err
	}

	app.Preferences().SetString(PREF_CERT, string(certPEM))
	app.Preferences().SetString(PREF_KEY, string(keyPEM))

	return tls.X509KeyPair(certPEM, keyPEM)
}

func generatePEMKeys() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"LocalMail P2P"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * 10 * time.Hour), // 10 лет

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip := ipnet.IP.To4(); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			}
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}

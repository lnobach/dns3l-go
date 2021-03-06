package common

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ConvertCertBundleToPEMStr(bundle []*x509.Certificate) (string, error) {
	var buf bytes.Buffer
	for _, cert := range bundle {
		err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func RSAPrivKeyToStr(privKey *rsa.PrivateKey) string {
	keyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	keyStr := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	return string(keyStr)
}

func RSAPrivKeyFromStr(privKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privKey))
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ParseCertificatePEM(certificate []byte) ([]*x509.Certificate, error) {
	result := make([]*x509.Certificate, 0)
	decodeTodo := certificate
	for {
		var block *pem.Block
		block, decodeTodo = pem.Decode(decodeTodo)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			result = append(result, cert)
		}
	}

	return result, nil

}

type PEMResource struct {
	PEMData     string
	ContentType string
	Domains     []string
}

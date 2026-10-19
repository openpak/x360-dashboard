package dashboard

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
)

// SignManifest wraps authored XML in the native XSigner envelope. The kernel
// verifies the OpenPak header signature; XAM independently checks payload SHA-1.
func SignManifest(payload, keyPEM []byte) ([]byte, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("signing key is not PEM")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		parsed, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	}
	key, ok := parsed.(*rsa.PrivateKey)
	if err != nil || !ok || key.N.BitLen() != 2048 || key.E != 65537 {
		return nil, fmt.Errorf("signing key must be RSA-2048 with exponent 65537")
	}
	if err := key.Validate(); err != nil {
		return nil, err
	}
	header := make([]byte, 0x130)
	copy(header, "YTGR")
	binary.LittleEndian.PutUint32(header[4:], 1)
	binary.LittleEndian.PutUint32(header[8:], uint32(len(header)))
	binary.LittleEndian.PutUint32(header[0x14:], uint32(len(payload)))
	digest := sha1.Sum(payload)
	copy(header[0x1c:], digest[:])
	digest = sha1.Sum(header[:0x30])
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, digest[:])
	if err != nil {
		return nil, err
	}
	copy(header[0x30:], signature)
	return append(header, payload...), nil
}

// HomeXML returns the original authored manifest before its signing envelope.
func HomeXML() ([]byte, error) { return content.ReadFile("dashhome.xml") }

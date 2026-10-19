package dashboard

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"testing"
)

func TestNativeSigningEnvelope(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	payload, err := HomeXML()
	if err != nil {
		t.Fatal(err)
	}
	body, err := SignManifest(payload, keyPEM)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != len(payload)+0x130 || string(body[:4]) != "YTGR" || binary.LittleEndian.Uint32(body[4:]) != 1 || binary.LittleEndian.Uint32(body[8:]) != 0x130 || int(binary.LittleEndian.Uint32(body[0x14:])) != len(payload) {
		t.Fatal("invalid native envelope")
	}
	hash := sha1.Sum(payload)
	if !bytes.Equal(hash[:], body[0x1c:0x30]) {
		t.Fatal("wrong payload hash")
	}
	hash = sha1.Sum(body[:0x30])
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA1, hash[:], body[0x30:0x130]); err != nil {
		t.Fatal(err)
	}
	body[0x14] ^= 1
	hash = sha1.Sum(body[:0x30])
	if rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA1, hash[:], body[0x30:0x130]) == nil {
		t.Fatal("modified header accepted")
	}
	if _, err := SignManifest(payload, []byte("invalid")); err == nil {
		t.Fatal("invalid key accepted")
	}
}

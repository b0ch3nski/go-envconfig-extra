package envconfigext

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
)

var exampleContent = []byte("example file content")

func TestFileContentFromPath(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := tmp.WriteAt(exampleContent, 0); err != nil {
		t.Fatal(err)
	}
	testFileContent(t, []byte(tmp.Name()), exampleContent)
}

func TestFileContentFromBase64(t *testing.T) {
	b64Content := make([]byte, base64.StdEncoding.EncodedLen(len(exampleContent)))
	base64.StdEncoding.Encode(b64Content, exampleContent)

	testFileContent(t, b64Content, exampleContent)
}

func TestFileContentFromPlaintext(t *testing.T) {
	testFileContent(t, exampleContent, exampleContent)
}

func testFileContent(t *testing.T, input, expected []byte) {
	var fc FileContent
	if err := fc.UnmarshalText(input); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(expected, fc) {
		t.Fatalf("Not equal: \nexpected: %v\nactual  : %v", expected, fc)
	}
}

func TestX509CertOK(t *testing.T) {
	sn := big.NewInt(123)
	cert := x509.Certificate{SerialNumber: sn}

	keyPub, keyPrv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, keyPub, keyPrv)
	if err != nil {
		t.Fatal(err)
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	var x509c X509Cert
	if err := x509c.UnmarshalText(certPem); err != nil {
		t.Fatal(err)
	}

	if x509c.SerialNumber.Cmp(sn) != 0 {
		t.Fatalf("Not equal: \nexpected: %v\nactual  : %v", sn, x509c.SerialNumber)
	}
}

func TestX509CertErr(t *testing.T) {
	cases := map[string][]byte{
		"NoPEM":  []byte(exampleContent),
		"BadPEM": pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte(exampleContent)}),
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			var x509c X509Cert
			if err := x509c.UnmarshalText(input); err == nil {
				t.Fatal("Error expected, got nil")
			}
		})
	}
}

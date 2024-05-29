package envconfigext

import (
	"crypto/x509"
	"encoding"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/wkhere/dtf"
)

var _ encoding.TextUnmarshaler = (*FileContent)(nil)

// FileContent holds the content of a file read either from filesystem path or string (plain or base64 encoded).
type FileContent []byte

// UnmarshalText implements `encoding.TextUnmarshaler` interface for reading file content from FS or string.
func (fc *FileContent) UnmarshalText(text []byte) error {
	// try #1: assume input is a filesystem path, try to read file
	if data, err := os.ReadFile(string(text)); err == nil {
		*fc = FileContent(data)
		return nil
	}

	// try #2: assume input is base64 encoded string, try to decode it
	*fc = make(FileContent, base64.StdEncoding.DecodedLen(len(text)))
	if _, err := base64.StdEncoding.Decode(*fc, text); err == nil {
		return nil
	}

	// fallback to plain text, return input without any changes
	*fc = text
	return nil
}

var (
	_ encoding.TextUnmarshaler = (*X509Cert)(nil)
	_ fmt.Stringer             = (*X509Cert)(nil)
)

// X509Cert is a `x509.Certificate` read either from filesystem path or string (plain or base64 encoded).
type X509Cert x509.Certificate

// UnmarshalText implements `encoding.TextUnmarshaler` interface for reading X509 certificate from FS or string.
func (c *X509Cert) UnmarshalText(text []byte) error {
	var fc FileContent
	if err := fc.UnmarshalText(text); err != nil {
		return err
	}

	pemData, _ := pem.Decode(fc)
	if pemData == nil || len(pemData.Bytes) == 0 {
		return fmt.Errorf("certificate: no PEM data found")
	}

	cert, err := x509.ParseCertificate(pemData.Bytes)
	if err != nil {
		return fmt.Errorf("certificate: error while parsing: %w", err)
	}

	*c = X509Cert(*cert)
	return nil
}

// String implements `fmt.Stringer` interface.
func (c *X509Cert) String() string {
	return fmt.Sprintf("Subject=%s | SAN=%v | Issuer=%s | Algorithm=%s | ValidFor=%s",
		c.Subject.CommonName, c.DNSNames, c.Issuer.CommonName, c.PublicKeyAlgorithm, dtf.Fmt(time.Until(c.NotAfter)))
}

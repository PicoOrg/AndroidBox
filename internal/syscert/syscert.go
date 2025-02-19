package syscert

import (
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/PicoOrg/AndroidBox/internal/util"
)

type SysCert interface {
}

type syscert struct {
	logger util.Logger
}

func NewSysCert(logger util.Logger) SysCert {
	return &syscert{
		logger: logger,
	}
}

func (instance *syscert) Install(certpath string) (err error) {
	content, err := os.ReadFile(certpath)
	if err != nil {
		instance.logger.Error("read cert error", util.Fields{"path": certpath, "error": err})
		return
	}
	p, _ := pem.Decode(content)
	if p.Type == "CERTIFICATE" {
	} else if _, err := x509.ParseCertificate(content); err != nil {
		content = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: content})
	} else {

	}

	return
}

func (instance *syscert) isDerOrPem(certpath string) (t string, err error) {
	content, err := os.ReadFile(certpath)
	if err != nil {
		instance.logger.Error("read cert error", util.Fields{"path": certpath, "error": err})
		return
	}
	p, _ := pem.Decode(content)
	if p.Type == "CERTIFICATE" {
	} else if _, err := x509.ParseCertificate(content); err != nil {
		content = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: content})
	} else {
	}
	return
}

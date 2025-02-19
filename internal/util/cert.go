package util

import (
	"crypto/x509"
	"os"
)

type Cert interface {
}

type cert struct {
	logger   Logger
	certpath string
	certtype string
}

func NewCert(certpath string, logger Logger) *cert {
	return &cert{
		logger:   logger,
		certpath: certpath,
	}
}

func (instance *cert) GetType() (t string, err error) {
	if instance.certtype != "" {
		return
	}
	content, err := os.ReadFile(instance.certpath)
	if err != nil {
		instance.logger.Error("read cert error", Fields{"path": instance.certpath, "error": err})
		return
	}

	if _, err = x509.ParseCertificate(content); err != nil {
		instance.logger.Error("read cert error", Fields{"path": instance.certpath, "error": err})
	}
	return
}

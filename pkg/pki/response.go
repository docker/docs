package pki

type CertificateResponse struct {
	Certificate      string `json:"certificate,omitempty"`
	CertificateChain string `json:"certificate_chain,omitempty"`
}

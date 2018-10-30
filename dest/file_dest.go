package dest

type FileDestType string

var (
	CAFile                FileDestType = "{{ .CA }}"
	CertificateFile       FileDestType = "{{ .Certificate }}"
	CertificateBundleFile FileDestType = "{{ .Certificate }}\n{{ .CABundle }}"
	PrivateKeyFile        FileDestType = "{{ .PrivateKey }}"
)

func (f FileDestType) Load() (string, error) {
	return string(f), nil
}

package dest

import (
	"strings"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type testDest struct {
	strings.Builder
}

func (d *testDest) Close() error {
	return nil
}

func (d *testDest) Write(p []byte) (n int, err error) {
	return d.Builder.Write(p)
}

func TestTemplateOutputs(t *testing.T) {
	tests := []struct {
		loader TemplateLoader
		res    *provider.Response
		expect string
	}{
		{
			TemplateString("{{ .CA }}"),
			&provider.Response{CA: "test-ca"},
			"test-ca",
		},
		{
			TemplateString("{{ .Certificate }}"),
			&provider.Response{Certificate: "test-cert"},
			"test-cert",
		},
		{
			TemplateString("{{ .CABundle }}"),
			&provider.Response{CABundle: "test-bundle"},
			"test-bundle",
		},
		{
			TemplateString("{{ .PrivateKey}}"),
			&provider.Response{PrivateKey: "private-key"},
			"private-key",
		},
		{
			TemplateString("{{ .Serial}}"),
			&provider.Response{Serial: "serial"},
			"serial",
		},
		{
			TemplateString("{{ .Certificate }}\n{{ .CABundle }}"),
			&provider.Response{CABundle: "ca-bundle", Certificate: "certificate"},
			"certificate\nca-bundle",
		},
	}

	for _, test := range tests {
		var d testDest
		tmpl := NewTemplate(test.loader, &d)
		err := tmpl.Send(test.res)
		if err != nil {
			t.Fatalf("unexpected error when using dest template: %v", err)
		}
		if d.Builder.String() != test.expect {
			t.Fatalf("dest template output \"%s\" was not expected \"%s\"",
				d.Builder.String(), test.expect)
		}
	}
}

package vagrant

import (
	"strings"
	"testing"
	"text/template"

	"github.com/TomHoenderdos/vbox/internal/config"
)

func renderTemplate(data vagrantfileData) (string, error) {
	tmpl, err := template.New("Vagrantfile").Parse(vagrantfileTmpl)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func TestVagrantfileTemplateBasic(t *testing.T) {
	data := vagrantfileData{
		Memory: 2048,
		CPUs:   2,
		Ports: []config.Port{
			{Guest: 4000, Host: 4000, Label: "Phoenix"},
		},
		Provision: "    apt-get update\n",
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	checks := []string{
		`prl.memory = 2048`,
		`prl.cpus = 2`,
		`guest: 4000, host: 4000`,
		`# Phoenix`,
		`apt-get update`,
		`echo "vbox provisioning complete!"`,
		`luminositylabsllc/bento-ubuntu-24.04-arm64`,
		`rsync__exclude`,
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("output missing %q", c)
		}
	}
}

func TestVagrantfileTemplateNoUSB(t *testing.T) {
	data := vagrantfileData{
		Memory:   2048,
		CPUs:     2,
		NeedsUSB: false,
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	if strings.Contains(out, "device-add") {
		t.Error("should not contain USB customization when NeedsUSB is false")
	}
}

func TestVagrantfileTemplateWithUSB(t *testing.T) {
	data := vagrantfileData{
		Memory:   2048,
		CPUs:     2,
		NeedsUSB: true,
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	if !strings.Contains(out, `"post-import"`) {
		t.Error("should contain USB customization when NeedsUSB is true")
	}
	if !strings.Contains(out, `--device-add`) {
		t.Error("should contain --device-add for USB")
	}
}

func TestVagrantfileTemplateMultiplePorts(t *testing.T) {
	data := vagrantfileData{
		Memory: 4096,
		CPUs:   4,
		Ports: []config.Port{
			{Guest: 4000, Host: 4000, Label: "Phoenix"},
			{Guest: 5432, Host: 15432, Label: "PostgreSQL"},
			{Guest: 8080, Host: 8080, Label: "HTTP"},
		},
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	if !strings.Contains(out, "guest: 4000, host: 4000") {
		t.Error("missing Phoenix port")
	}
	if !strings.Contains(out, "guest: 5432, host: 15432") {
		t.Error("missing PostgreSQL port")
	}
	if !strings.Contains(out, "guest: 8080, host: 8080") {
		t.Error("missing HTTP port")
	}
	if !strings.Contains(out, "prl.memory = 4096") {
		t.Error("wrong memory")
	}
	if !strings.Contains(out, "prl.cpus = 4") {
		t.Error("wrong cpus")
	}
}

func TestVagrantfileTemplateNoPorts(t *testing.T) {
	data := vagrantfileData{
		Memory: 2048,
		CPUs:   2,
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	if strings.Contains(out, "forwarded_port") {
		t.Error("should not contain forwarded_port when no ports configured")
	}
}

func TestVagrantfileTemplateClaudeConfig(t *testing.T) {
	data := vagrantfileData{
		Memory: 2048,
		CPUs:   2,
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	checks := []string{
		`.claude"`,
		`.claude.json`,
		`.config/gh`,
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("output missing Claude/tool config sync: %q", c)
		}
	}
}

func TestVagrantfileTemplateProvision(t *testing.T) {
	provision := `    apt-get update
    apt-get install -y golang
    su - vagrant -c 'go version'
`
	data := vagrantfileData{
		Memory:    2048,
		CPUs:      2,
		Provision: provision,
	}

	out, err := renderTemplate(data)
	if err != nil {
		t.Fatalf("template error: %v", err)
	}

	if !strings.Contains(out, "apt-get install -y golang") {
		t.Error("provision script not included in output")
	}
}

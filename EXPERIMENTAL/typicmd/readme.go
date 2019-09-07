package typicmd

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typirecipe/readme"
)

const (
	readmeFile     = "README.md"
	configTemplate = `| Key | Type | Default | Required | Description |	
|---|---|---|---|---|{{range .}}
|{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_required .}}|{{usage_description .}}|{{end}}`
)

// GenerateReadme for generate typical applical readme
func generateReadme(a *typictx.ActionContext) (err error) {
	readme0 := readme.DefaultReadme().
		SetTitle(a.Name).
		SetDescription(a.Description).
		SetSection("Configuration", func(md *typirecipe.Markdown) (err error) {
			for _, acc := range a.Context.ConfigAccessors() {
				name := acc.GetName()
				if name != "" {
					md.Heading3(name)
				}
				var builder strings.Builder
				envconfig.Usagef(acc.GetConfigPrefix(), acc.GetConfigSpec(), &builder, configTemplate)
				md.Writeln(builder.String())
			}
			return
		})
	log.Infof("Generate new %s", readmeFile)
	err = readme0.OutputToFile(readmeFile)
	if err != nil {
		return
	}
	return
}

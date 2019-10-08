package prebuilder

import (
	log "github.com/sirupsen/logrus"

	"time"

	"github.com/typical-go/runn"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/bash"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/golang"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typicmd/internal/prebuilder/walker"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typictx"
	"github.com/typical-go/typical-rest-server/EXPERIMENTAL/typienv"
)

// PreBuilder responsible to prebuild process
type PreBuilder struct {
	*typictx.Context
	*walker.ProjectFiles
	*walker.ContextFile
	Filenames []string
	Packages  []string
}

// TestTargets generate test target
func (p *PreBuilder) TestTargets() error {
	defer elapsed("Generate TestTargets")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	src.AddTestTargets(p.Packages...)
	target := dependency + "/test_targets.go"
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

// Annotated to generate annotated
func (p *PreBuilder) Annotated() error {
	defer elapsed("Generate Annotated")()
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg)
	for _, pkg := range p.Packages {
		src.AddImport("", p.Context.Root+"/"+pkg)
	}
	src.AddConstructors(p.ProjectFiles.Autowires()...)
	src.AddMockTargets(p.ProjectFiles.Automocks()...)
	target := dependency + "/annotateds.go"
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

// Configuration to generate configuration
func (p *PreBuilder) Configuration() error {
	defer elapsed("Generate Configuration")()
	conf := createConfiguration(p.Context)
	pkg := typienv.Dependency.Package
	src := golang.NewSourceCode(pkg).AddStruct(conf.Struct)
	src.AddImport("", "github.com/kelseyhightower/envconfig")
	for _, imp := range p.ContextFile.Imports {
		src.AddImport(imp.Name, imp.Path)
	}
	src.AddConstructors(conf.Constructors...)
	target := dependency + "/configurations.go"
	return runn.Execute(
		src.Cook(target),
		bash.GoImports(target),
	)
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Infof("%s took %v\n", what, time.Since(start))
	}
}
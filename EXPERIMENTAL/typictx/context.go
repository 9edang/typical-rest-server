package typictx

import (
	"os"
	"strings"

	"go.uber.org/dig"
	"gopkg.in/urfave/cli.v1"
)

// Context of typical application
type Context struct {
	Name        string
	Version     string
	Description string
	AppModule   AppModule

	ReadmeTemplate string
	ReadmeFile     string

	AppPkg     string
	MockPkg    string
	BinaryName string

	Modules []*Module

	Commands []cli.Command
}

// ReadmeTemplateOrDefault return readme template field or the default template
func (c Context) ReadmeTemplateOrDefault() string {
	if c.ReadmeTemplate == "" {
		return defaultReadmeTemplate
	}
	return c.ReadmeTemplate
}

// ReadmeFileOrDefault return readme file field or default template
func (c Context) ReadmeFileOrDefault() string {
	if c.ReadmeFile == "" {
		return defaultReadmeFile
	}

	return c.ReadmeFile
}

// BinaryNameOrDefault return binary name of typiapp or default value
func (c Context) BinaryNameOrDefault() string {
	if c.BinaryName == "" {
		dir, _ := os.Getwd()
		chunks := strings.Split(dir, "/")
		return chunks[len(chunks)-1]
	}
	return c.BinaryName
}

// AppPkgOrDefault return application package of typiapp or default value
func (c Context) AppPkgOrDefault() string {
	if c.AppPkg == "" {
		return defaultApplicationPkg
	}
	return c.AppPkg
}

// MockPkgOrDefault return mock package of typiapp or default value
func (c Context) MockPkgOrDefault() string {
	if c.MockPkg == "" {
		return defaultMockPkg
	}
	return c.MockPkg
}

// Container to return the depedency injection
func (c Context) Container() *dig.Container {
	container := dig.New()

	for _, constructor := range c.AppModule.GetConstructors() {
		container.Provide(constructor)
	}

	for _, module := range c.Modules {
		module.Inject(container)
	}

	return container
}
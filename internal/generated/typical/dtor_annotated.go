package typical

/* Autogenerated by Typical-Go. DO NOT EDIT.

TagName:
	@dtor

Help:
	https://pkg.go.dev/github.com/typical-go/typical-go/pkg/typapp?tab=doc#DtorAnnotation
*/

import (
	"github.com/typical-go/typical-go/pkg/typapp"
	a "github.com/typical-go/typical-rest-server/internal/app/infra"
)

func init() {
	typapp.AppendDtor(
		&typapp.Destructor{Fn: a.Teardown},
	)
}

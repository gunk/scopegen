package generate

import (
	"fmt"
	"html/template"
	"io"

	"github.com/gunk/scopegen/parser"
)

var (
	goTmplV1 = template.Must(template.New("").Parse(`// Code generated by "scopegen"; DO NOT EDIT.
package {{ .Package }}
type AuthScope string
type ServiceScope struct{}
{{ if .Scopes -}}
const (
{{- range $scope := .Scopes }}
	Scope_{{ $scope.Name }} AuthScope = "{{ $scope.Value }}"
{{- end }}
)
{{ end -}}
var authScopes = map[string][]AuthScope{
{{- range $method := .Methods }}
	"{{ $method.Name }}": { {{- range $index, $scope := $method.Scopes }}{{ if $index }}, {{ end }}Scope_{{ $scope }}{{ end -}} },
{{- end }}
}
// Any allows a loose challenge, for claims containing any of the method scopes.
//
// Use All method as a default for OAuth2 style scopes.  Any is useful with more complex scope definitions.
func (svcSc ServiceScope) Any(method string, claims []string) bool {
	ch := authScopes[method]
	for _, s := range ch {
		for _, c := range claims {
			if string(s) == c {
				return true
			}
		}
	}
	return len(ch) == 0
}
// All is the default OAuth2 challenge method, ensuring claims contains all method scopes.
func (svcSc ServiceScope) All(method string, claims []string) bool {
	ch := authScopes[method]
	if len(ch) > len(claims) {
		return false
	}
scopes:
	for _, s := range ch {
		for _, c := range claims {
			if string(s) == c {
				continue scopes
			}
		}
		return false
	}
	return true
}
`))
	goTmplV2 = template.Must(template.New("").Parse(`// Code generated by "scopegen"; DO NOT EDIT.
package {{ .Package }}
type ServiceScope struct{}
var Scopes = map[string]string{
{{- range $scope := .Scopes }}
	"{{ $scope.Name }}": "{{ $scope.Value }}",
{{- end }}
}
var AuthScopes = map[string][]string{
{{- range $method := .Methods }}
	"{{ $method.Name }}": { {{- range $index, $scope := $method.Scopes }}{{ if $index }}, {{ end }}"{{ $scope }}"{{ end -}} },
{{- end }}
}
// Any allows a loose challenge, for claims containing any of the method scopes.
//
// Use All method as a default for OAuth2 style scopes.  Any is useful with more complex scope definitions.
func (svcSc *ServiceScope) Any(method string, claims []string) bool {
	ch := AuthScopes[method]
	for _, s := range ch {
		for _, c := range claims {
			if string(s) == c {
				return true
			}
		}
	}
	return len(ch) == 0
}
// All is the default OAuth2 challenge method, ensuring claims contains all method scopes.
func (svcSc *ServiceScope) All(method string, claims []string) bool {
	ch := AuthScopes[method]
	if len(ch) > len(claims) {
		return false
	}
scopes:
	for _, s := range ch {
		for _, c := range claims {
			if string(s) == c {
				continue scopes
			}
		}
		return false
	}
	return true
}
`))
)

// GO generates defined OAuth2 scopes for go programming languages.
func Go(w io.Writer, f *parser.File, outputVersion int) error {
	switch outputVersion {
	case 1:
		return goTmplV1.Execute(w, f)
	case 2:
		return goTmplV2.Execute(w, f)
	default:
		return fmt.Errorf("unknown ouput version: %d", outputVersion)
	}
}

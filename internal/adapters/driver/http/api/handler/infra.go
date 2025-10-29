package handler

import (
	"marketflow/internal/ports/inbound"
	"text/template"
)

// import (
// 	"text/template"

// 	"marketflow/internal/ports/inbound"
// )

type handler struct {
	use inbound.UsecaseInter
}

func NewHandler(use inbound.UsecaseInter) any {
	templates, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		return nil
	}
	// соңғысының аты болады
	// fmt.Println(templates.Name())
	// for _, t := range templates.Templates() {
	// 	fmt.Println(t.Name())
	// }
	return &handler{templates}
}

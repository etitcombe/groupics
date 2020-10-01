package main

import (
	"html/template"
	"time"

	"github.com/etitcombe/groupics/pkg/forms"
	"github.com/etitcombe/groupics/pkg/models"
)

type viewModel struct {
	CSRFToken       string
	Flash           template.HTML
	IsAuthenticated bool
	Year            int
}

type createViewModel struct {
	viewModel
	Form *forms.Form
}

type homeViewModel struct {
	viewModel
	Snippets []*models.Snippet
}

type showViewModel struct {
	viewModel
	Snippet *models.Snippet
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

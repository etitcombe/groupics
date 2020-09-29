package main

import (
	"html/template"

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

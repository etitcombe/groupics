package main

import (
	"net/url"

	"github.com/etitcombe/groupics/pkg/models"
)

type viewModel struct {
	Year int
}

type createViewModel struct {
	viewModel
	FormData   url.Values
	FormErrors map[string]string
}

type homeViewModel struct {
	viewModel
	Snippets []*models.Snippet
}

type showViewModel struct {
	viewModel
	Snippet *models.Snippet
}

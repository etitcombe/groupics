package main

import "github.com/etitcombe/groupics/pkg/models"

type viewModel struct {
	Year int
}
type homeViewModel struct {
	viewModel
	Snippets []*models.Snippet
}

type showViewModel struct {
	viewModel
	Snippet *models.Snippet
}

package main

import "github.com/etitcombe/groupics/pkg/models"

type homeViewModel struct {
	Snippets []*models.Snippet
}

type showViewModel struct {
	Snippet *models.Snippet
}

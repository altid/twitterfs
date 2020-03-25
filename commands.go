package main

import "github.com/altid/libs/fs"

var TwitterCommands = []*fs.Command{
	{
		Name:        "tweet",
		Args:        []string{},
		Heading:     fs.DefaultGroup,
		Description: "Send a Tweet",
	},
	{
		Name:        "rt",
		Args:        []string{"<#id>"},
		Heading:     fs.DefaultGroup,
		Description: "Retweet by ID",
	},
}
package app

import "html/template"

type BlogPost struct {
	Title                   string
	Slug                    string
	Parent                  string
	Content                 template.HTML
	Description             string
	Order                   int
	Headers                 []string
	MetaDescription         string
	MetaPropertyTitle       string
	MetaPropertyDescription string
	MetaOgURL               string
}

type SideBarData struct {
	Categories []Category
}

type Category struct {
	Name  string
	Pages []BlogPost
	Order int
}

package posts

import "embed"

//go:embed *.md

var Assets embed.FS

file, err := Assets.ReadFile("my-post.md")
if err != nil {
  return err
}

return string(file)

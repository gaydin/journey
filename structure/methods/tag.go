package methods

import (
	"context"
	"strings"

	"github.com/gaydin/journey/slug"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

func GenerateTagsFromCommaString(input string) []structure.Tag {
	output := make([]structure.Tag, 0)
	tags := strings.Split(input, ",")
	for index, _ := range tags {
		tags[index] = strings.TrimSpace(tags[index])
	}
	for _, tag := range tags {
		if tag != "" {
			output = append(output, structure.Tag{Name: tag, Slug: slug.Generate(context.Background(), store.DB, tag, "tags")})
		}
	}
	return output
}

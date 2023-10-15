package methods

import (
	"context"
	"strings"

	"github.com/gaydin/journey/slug"
	"github.com/gaydin/journey/store"
	"github.com/gaydin/journey/structure"
)

func GenerateTagsFromCommaString(ctx context.Context, db store.Database, input string) []structure.Tag {
	output := make([]structure.Tag, 0)
	tags := strings.Split(input, ",")
	for index := range tags {
		tags[index] = strings.TrimSpace(tags[index])
	}
	for _, tag := range tags {
		if tag != "" {
			output = append(output, structure.Tag{Name: tag, Slug: slug.Generate(ctx, db, tag, "tags")})
		}
	}

	return output
}

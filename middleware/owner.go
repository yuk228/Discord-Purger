package middleware

import (
	"os"
	"slices"
	"strings"

	"github.com/yuk228/disgolf"
)

func HasOwnerMiddleware(ctx *disgolf.MessageCtx) {
	if slices.Contains(strings.Split(os.Getenv("OWNER_IDS"), ","), ctx.Message.Author.ID) {
		ctx.Next()
	}
}


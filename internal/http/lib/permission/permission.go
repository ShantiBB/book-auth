package permission

import (
	"context"
)

func Admin(ctx context.Context) bool {
	if role := ctx.Value("userRole").(string); role != "admin" {
		return false
	}
	return true
}

func Moderator(ctx context.Context) bool {
	if role := ctx.Value("userRole").(string); role != "moderator" {
		return false
	}
	return true
}

func UserOwn(ctx context.Context, userID int64) bool {
	if id := ctx.Value("userID").(int64); id != userID {
		return false
	}
	return true
}

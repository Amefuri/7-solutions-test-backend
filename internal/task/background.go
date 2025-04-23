package task

import (
	"7-solutions-test-backend/internal/core/user"
	"context"
	"fmt"
	"time"
)

func StartUserCountLogger(ctx context.Context, repo user.Repository) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			count, _ := repo.Count(ctx)
			fmt.Println("📊 Current user count:", count)
		case <-ctx.Done():
			return
		}
	}
}

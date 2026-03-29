package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"golang_world_cup/internal/user"
)

// InitScheduler 啟動每日積分發放排程
func InitScheduler(db *gorm.DB) {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local),
	)

	_, err := c.AddFunc("0 0 0 * * *", func() {
		log.Println("[Scheduler] Running daily points distribution...")
		distributeDailyPoints(db)
	})

	if err != nil {
		log.Fatalf("[Scheduler] Failed to add cron job: %v", err)
	}

	c.Start()
	log.Println("[Scheduler] Started successfully.")
}

func distributeDailyPoints(db *gorm.DB) {
	result := db.Model(&user.User{}).Update("points", gorm.Expr("points + ?", 1000))
	if result.Error != nil {
		log.Printf("[Scheduler] Error updating daily points: %v\n", result.Error)
	} else {
		log.Printf("[Scheduler] Successfully added 1000 points to %d users.\n", result.RowsAffected)
	}
}

package worker

import (
	"chat-golang/src/internal/domain/entities"
	"chat-golang/src/internal/infrastructure/database"
	"log"
	"time"
)

func StartReminderWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			checkAgendas()
		}
	}()
	log.Println("Reminder worker started")
}

func checkAgendas() {
	now := time.Now()
	upcoming := now.Add(15 * time.Minute)

	var agendas []entities.Agenda
	// Find agendas that are scheduled within the next 15 minutes and haven't been notified (we should add a notified flag later)
	err := database.DB.Where("scheduled_at > ? AND scheduled_at <= ?", now, upcoming).Find(&agendas).Error
	if err != nil {
		log.Printf("Worker error fetching agendas: %v", err)
		return
	}

	for _, agenda := range agendas {
		log.Printf("REMINDER: Agenda '%s' is starting at %v in room %s", agenda.Title, agenda.ScheduledAt, agenda.RoomID)
		// TODO: Send WebSocket notification to room members
	}
}

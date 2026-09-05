package service

import (
	"context"
	"time"

	"github.com/MamangRust/microservice-ecommerce-pkg/database/models"
	"github.com/MamangRust/microservice-ecommerce-pkg/outbox"
	"gorm.io/gorm"
)

// outboxQuerier adapts GORM to the outbox.OutboxQuerier contract.
type outboxQuerier struct {
	db *gorm.DB
}

func NewOutboxQuerier(db *gorm.DB) outbox.OutboxQuerier {
	if db == nil {
		return nil
	}
	return &outboxQuerier{db: db}
}

func (q *outboxQuerier) CreateOutboxEvent(ctx context.Context, params outbox.CreateOutboxEventParams) (outbox.OutboxEvent, error) {
	event := models.OutboxEvent{
		Topic:    params.Topic,
		EventKey: params.EventKey,
		Payload:  params.Payload,
		Status:   "pending",
	}
	if err := q.db.WithContext(ctx).Create(&event).Error; err != nil {
		return outbox.OutboxEvent{}, err
	}
	return outbox.OutboxEvent{
		OutboxID: event.OutboxID,
		Topic:    event.Topic,
		EventKey: event.EventKey,
		Payload:  event.Payload,
		Status:   event.Status,
	}, nil
}

func (q *outboxQuerier) ClaimPendingOutboxEvents(ctx context.Context, params outbox.ClaimPendingOutboxEventsParams) ([]outbox.OutboxEvent, error) {
	var events []models.OutboxEvent
	err := q.db.WithContext(ctx).Where("status = 'pending' AND (next_attempt_at IS NULL OR next_attempt_at <= ?)", params.NextAttemptAt).
		Order("outbox_id ASC").Limit(int(params.Limit)).Find(&events).Error
	if err != nil {
		return nil, err
	}

	result := make([]outbox.OutboxEvent, 0, len(events))
	for _, e := range events {
		result = append(result, outbox.OutboxEvent{
			OutboxID: e.OutboxID,
			Topic:    e.Topic,
			EventKey: e.EventKey,
			Payload:  e.Payload,
			Status:   e.Status,
		})
	}
	return result, nil
}

func (q *outboxQuerier) MarkOutboxEventFailed(ctx context.Context, params outbox.MarkOutboxEventFailedParams) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).First(&event, params.OutboxID).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	event.Status = "failed"
	event.Attempts++
	na := params.NextAttemptAt
	event.NextAttemptAt = na
	q.db.WithContext(ctx).Save(&event)
	return outbox.OutboxEvent{
		OutboxID: event.OutboxID,
		Topic:    event.Topic,
		EventKey: event.EventKey,
		Payload:  event.Payload,
		Status:   event.Status,
	}, nil
}

func (q *outboxQuerier) MarkOutboxEventDelivered(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).First(&event, outboxID).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	event.Status = "delivered"
	q.db.WithContext(ctx).Save(&event)
	return outbox.OutboxEvent{
		OutboxID: event.OutboxID,
		Topic:    event.Topic,
		EventKey: event.EventKey,
		Payload:  event.Payload,
		Status:   event.Status,
	}, nil
}

func (q *outboxQuerier) MarkOutboxEventDead(ctx context.Context, outboxID int64) (outbox.OutboxEvent, error) {
	var event models.OutboxEvent
	err := q.db.WithContext(ctx).First(&event, outboxID).Error
	if err != nil {
		return outbox.OutboxEvent{}, err
	}
	event.Status = "dead"
	q.db.WithContext(ctx).Save(&event)
	return outbox.OutboxEvent{
		OutboxID: event.OutboxID,
		Topic:    event.Topic,
		EventKey: event.EventKey,
		Payload:  event.Payload,
		Status:   event.Status,
	}, nil
}

func (q *outboxQuerier) DeleteOldOutboxEvents(ctx context.Context, cutoff time.Time) (int64, error) {
	result := q.db.WithContext(ctx).Where("created_at < ? AND status = 'delivered'", cutoff).Delete(&models.OutboxEvent{})
	return result.RowsAffected, result.Error
}

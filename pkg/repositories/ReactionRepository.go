package repositories

import (
	"database/sql"
	realtimeforum "real-time-forum"
	"real-time-forum/pkg/models"
)

type ReactionRepository interface {
	UpsertReaction(userId string, entityType string, entityId int, score int) (int, error)
	DoesCommentExists(commentId int) error
}

func (db *DB) UpsertReaction(userId string, entityType string, entityId int, score int) (int, error) {
	switch entityType {
	case "post":
		if err := db.DoesPostExists(entityId); err != nil {
			return 0, err
		}
	case "comment":
		if err := db.DoesCommentExists(entityId); err != nil {
			return 0, err
		}
	default:
		return 0, realtimeforum.ErrBadRequest
	}

	var existingReaction models.Reaction
	err := db.Conn.QueryRow(
		`SELECT reactionId, userId, entityType, entityId, score, createdAt
		 FROM reaction
		 WHERE userId = ? AND entityType = ? AND entityId = ?`,
		userId, entityType, entityId,
	).Scan(
		&existingReaction.ReactionId,
		&existingReaction.UserId,
		&existingReaction.EntityType,
		&existingReaction.EntityId,
		&existingReaction.Score,
		&existingReaction.CreatedAt,
	)

	if err == sql.ErrNoRows {
		_, err := db.Conn.Exec(
			`INSERT INTO reaction (userId, entityType, entityId, score, createdAt)
			 VALUES (?, ?, ?, ?, datetime('now'))`,
			userId, entityType, entityId, score,
		)
		if err != nil {
			return 0, realtimeforum.ErrInternal
		}
	} else if err != nil {
		return 0, realtimeforum.ErrInternal
	} else {
		if existingReaction.Score == score {
			_, err := db.Conn.Exec(
				`DELETE FROM reaction WHERE reactionId = ?`,
				existingReaction.ReactionId,
			)
			if err != nil {
				return 0, realtimeforum.ErrInternal
			}
			score = 0
		} else {
			_, err := db.Conn.Exec(
				`UPDATE reaction SET score = ?, createdAt = datetime('now') WHERE reactionId = ?`,
				score, existingReaction.ReactionId,
			)
			if err != nil {
				return 0, realtimeforum.ErrInternal
			}
		}
	}

	var totalScore int
	err = db.Conn.QueryRow(
		`SELECT COALESCE(SUM(score), 0) FROM reaction WHERE entityType = ? AND entityId = ?`,
		entityType, entityId,
	).Scan(&totalScore)
	if err != nil {
		return 0, realtimeforum.ErrInternal
	}

	return totalScore, nil
}

func (db *DB) DoesCommentExists(commentId int) error {
	var count int
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM comment WHERE commentId = ?", commentId).Scan(&count)
	if err != nil {
		return realtimeforum.ErrInternal
	}
	if count == 0 {
		return realtimeforum.ErrNotFound
	}
	return nil
}
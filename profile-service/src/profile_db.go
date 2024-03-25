package src

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Cprime50/user/db"
	pb "github.com/Cprime50/user/profilepb"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrProfileNotFound           = errors.New("profile not found")
	ErrDuplicateEntry            = errors.New("duplicate entry")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)

func getProfileByUserId(userId string) (*pb.Profile, error) {
	profile := pb.Profile{}
	var createdAt, updatedAt time.Time
	err := db.Db.QueryRow("SELECT id, user_id, email, username, bio, avatar, score, created_at, updated_at FROM profiles WHERE user_id = $1", userId).
		Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.Username, &profile.Bio, &profile.Avatar, &profile.Score, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProfileNotFound
		}
		return nil, fmt.Errorf("getProfile: %w", err)
	}

	// Convert time.Time to timestamppb.Timestamp
	profile.CreatedAt = timestamppb.New(createdAt)
	profile.UpdatedAt = timestamppb.New(updatedAt)

	return &profile, nil

}

func updateProfile(p *pb.Profile) error {
	result, err := db.Db.Exec(
		"UPDATE profiles SET username = $1, bio = $2, avatar = $3, updated_at = $4 WHERE user_id = $5",
		p.Username,
		p.Bio,
		p.Avatar,
		time.Now(),
		p.UserId,
	)
	if err != nil {
		return fmt.Errorf("UpdateProfile error: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrProfileNotFound
	}
	return nil
}

func createProfile(p *pb.Profile) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("uuid.NewRandom: %w", err)
	}
	_, err = db.Db.Exec(
		"INSERT INTO profiles (id, user_id, email, username, avatar, bio, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		id,
		p.UserId,
		p.Email,
		p.Username,
		p.Avatar,
		p.Bio,
		time.Now(),
	)
	if err != nil {
		sqliteErr, _ := err.(sqlite3.Error)
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return ErrDuplicateEntry
		}
		return fmt.Errorf("CreateProfile error: %w", err)
	}
	return nil
}

func selectProfiles() ([]*pb.Profile, error) {
	rows, err := db.Db.Query("SELECT * FROM profiles")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	var createdAt, updatedAt time.Time
	var profiles []*pb.Profile
	for rows.Next() {
		profile := &pb.Profile{}

		if err := rows.Scan(&profile.Id, &profile.UserId, &profile.Email, &profile.Username, &profile.Bio, &profile.Avatar, &profile.Score, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		// Convert time.Time to timestamppb.Timestamp
		profile.CreatedAt = timestamppb.New(createdAt)
		profile.UpdatedAt = timestamppb.New(updatedAt)
		profiles = append(profiles, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}

	if len(profiles) == 0 {
		return nil, ErrProfileNotFound
	}

	return profiles, nil
}

func deleteProfileByUserId(UserID string) error {
	result, err := db.Db.Exec("DELETE FROM profiles WHERE user_id = $1", UserID)
	if err != nil {
		return fmt.Errorf("error deleting profile: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrProfileNotFound
	}
	return nil
}

func updateScore(userId string, score int64) error {
	result, err := db.Db.Exec("UPDATE profiles SET score = $1 WHERE user_id = $2", score, userId)
	if err != nil {
		return fmt.Errorf("error updating score: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrProfileNotFound
	}
	return nil
}

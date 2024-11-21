package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserSecurityQuestionRepository struct{}

func NewUserSecurityQuestionRepository() *UserSecurityQuestionRepository {
	return &UserSecurityQuestionRepository{}
}

func (r *UserSecurityQuestionRepository) Insert(ctx context.Context, userSecurityQuestion *models.UserSecurityQuestion) (models.UserSecurityQuestion, error) {
	query := `
		INSERT INTO auth.user_security_questions (user_id, question, answer, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, userSecurityQuestion.UserID, userSecurityQuestion.Question, userSecurityQuestion.Answer, time.Now(), time.Now()).Scan(&userSecurityQuestion.ID)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("failed to insert user security question: %w", err)
	}

	return *userSecurityQuestion, nil
}

func (r *UserSecurityQuestionRepository) Query(ctx context.Context, id uuid.UUID) (models.UserSecurityQuestion, error) {
	query := `
		SELECT id, user_id, question, answer, created_at, updated_at
		FROM auth.user_security_questions WHERE id = $1`
	pool := db.GetPools().ReadPool

	userSecurityQuestion := models.UserSecurityQuestion{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&userSecurityQuestion.ID,
		&userSecurityQuestion.UserID,
		&userSecurityQuestion.Question,
		&userSecurityQuestion.Answer,
		&userSecurityQuestion.CreatedAt,
		&userSecurityQuestion.UpdatedAt,
	)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("failed to query user security question: %w", err)
	}

	return userSecurityQuestion, nil
}

func (r *UserSecurityQuestionRepository) Update(ctx context.Context, userSecurityQuestion *models.UserSecurityQuestion) (models.UserSecurityQuestion, error) {
	query := `
		UPDATE auth.user_security_questions
		SET user_id = $1, question = $2, answer = $3, updated_at = $4
		WHERE id = $5 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedUserSecurityQuestion := models.UserSecurityQuestion{}
	err := pool.QueryRow(ctx, query, userSecurityQuestion.UserID, userSecurityQuestion.Question, userSecurityQuestion.Answer, time.Now(), userSecurityQuestion.ID).Scan(&updatedUserSecurityQuestion.ID)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("failed to update user security question: %w", err)
	}

	return updatedUserSecurityQuestion, nil
}

func (r *UserSecurityQuestionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.user_security_questions WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user security question: %w", err)
	}

	return nil
}

func (r *UserSecurityQuestionRepository) QueryByUser(ctx context.Context, userID uuid.UUID) ([]models.UserSecurityQuestion, error) {
	query := `SELECT id, user_id, question, answer, created_at, updated_at FROM auth.user_security_questions WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user security questions by user ID: %w", err)
	}
	defer rows.Close()

	var userSecurityQuestions []models.UserSecurityQuestion
	for rows.Next() {
		var userSecurityQuestion models.UserSecurityQuestion
		err := rows.Scan(
			&userSecurityQuestion.ID,
			&userSecurityQuestion.UserID,
			&userSecurityQuestion.Question,
			&userSecurityQuestion.Answer,
			&userSecurityQuestion.CreatedAt,
			&userSecurityQuestion.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user security question row: %w", err)
		}
		userSecurityQuestions = append(userSecurityQuestions, userSecurityQuestion)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user security question rows: %w", err)
	}

	return userSecurityQuestions, nil
}

func (r *UserSecurityQuestionRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM auth.user_security_questions WHERE user_id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user security questions by user ID: %w", err)
	}

	return nil
}
func (r *UserSecurityQuestionRepository) GenerateSecurityQuestions(ctx context.Context, userID uuid.UUID, count int) ([]models.UserSecurityQuestion, error) {
	questions := []string{
		"What is your mother's maiden name?",
		"What is the name of your first pet?",
		"What is your favorite color?",
		"What is the name of your childhood best friend?",
		"What is your favorite food?",
		"What is the name of your high school?",
		"What is your favorite book?",
		"What is your favorite movie?",
		"What is your favorite song?",
		"What is your favorite hobby?",
	}

	if count > len(questions) {
		return nil, fmt.Errorf("cannot generate more security questions than available")
	}

	securityQuestions := make([]models.UserSecurityQuestion, count)
	for i := 0; i < count; i++ {
		securityQuestion := models.UserSecurityQuestion{
			UserID:   userID,
			Question: questions[i],
			Answer:   "", // Answer will be set by the user
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		insertedSecurityQuestion, err := r.Insert(ctx, &securityQuestion)
		if err != nil {
			return nil, fmt.Errorf("failed to insert security question: %w", err)
		}
		securityQuestions[i] = insertedSecurityQuestion
	}

	return securityQuestions, nil
}

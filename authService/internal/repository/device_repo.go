package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DeviceRepository struct{}

func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{}
}

func (r *DeviceRepository) Insert(ctx context.Context, device *models.Device) ( uuid.UUID, error) {
	query := `
		INSERT INTO auth.devices (user_id, os, browser, device_type, ip_address, user_agent, last_used_at, region)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, device.UserID, device.OS, device.Browser, device.DeviceType, device.IPAddress, device.UserAgent, device.LastUsedAt, device.Region).Scan(&device.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert device: %w", err)
	}

	return device.ID, nil
}

func (r *DeviceRepository) Query(ctx context.Context, id uuid.UUID) (models.Device, error) {
	query := `
		SELECT id, user_id, os, browser, device_type, ip_address, user_agent, last_used_at, region
		FROM auth.devices WHERE id = $1`
	pool := db.GetPools().ReadPool

	device := models.Device{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&device.ID,
		&device.UserID,
		&device.OS,
		&device.Browser,
		&device.DeviceType,
		&device.IPAddress,
		&device.UserAgent,
		&device.LastUsedAt,
		&device.Region,
	)
	if err != nil {
		return models.Device{}, fmt.Errorf("failed to query device: %w", err)
	}

	return device, nil
}

func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) (models.Device, error) {
	query := `
		UPDATE auth.devices
		SET user_id = $1, os = $2, browser = $3, device_type = $4, ip_address = $5, user_agent = $6, last_used_at = $7, region = $8
		WHERE id = $9 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedDevice := models.Device{}
	err := pool.QueryRow(ctx, query, device.UserID, device.OS, device.Browser, device.DeviceType, device.IPAddress, device.UserAgent, device.LastUsedAt, device.Region, device.ID).Scan(&updatedDevice.ID)
	if err != nil {
		return models.Device{}, fmt.Errorf("failed to update device: %w", err)
	}

	return updatedDevice, nil
}

func (r *DeviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.devices WHERE id = $1 RETURNING id`
	pool := db.GetPools().DeletePool

	var deletedDeviceID uuid.UUID
	err := pool.QueryRow(ctx, query, id).Scan(&deletedDeviceID)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil
}

func (r *DeviceRepository) QueryByUser(ctx context.Context, userID uuid.UUID) ([]models.Device, error) {
	query := `SELECT id, user_id, os, browser, device_type, ip_address, user_agent, last_used_at, region FROM auth.devices WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices by user ID: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID,
			&device.UserID,			&device.OS,
			&device.Browser,
			&device.DeviceType,
			&device.IPAddress,
			&device.UserAgent,
			&device.LastUsedAt,
			&device.Region,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device row: %w", err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over device rows: %w", err)
	}

	return devices, nil
}

func (r *DeviceRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM auth.devices WHERE user_id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete devices by user ID: %w", err)
	}

	return nil
}
func (r *DeviceRepository) UpdateLastUsedAt(ctx context.Context, id uuid.UUID, lastUsedAt time.Time) error {
	query := `UPDATE auth.devices SET last_used_at = $1 WHERE id = $2`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, lastUsedAt, id)
	if err != nil {
		return fmt.Errorf("failed to update device last used at: %w", err)
	}

	return nil
}
 
func (r *DeviceRepository) CountDevicesByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM auth.devices WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	var count int
	err := pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count devices by user ID: %w", err)
	}

	return count, nil
}
func (r *DeviceRepository) DeleteExpiredDevices(ctx context.Context, expirationTime time.Time) error {
	query := `DELETE FROM auth.devices WHERE last_used_at < $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, expirationTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired devices: %w", err)
	}

	return nil
}

func (r *DeviceRepository) QueryByIPAddress(ctx context.Context, ipAddress string) ([]models.Device, error) {
	query := `SELECT id, user_id, os, browser, device_type, ip_address, user_agent, last_used_at, region FROM auth.devices WHERE ip_address = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices by IP address: %w", err)
	}
	defer rows.Close()

	var devices []models.Device
	for rows.Next() {
		var device models.Device
		err := rows.Scan(
			&device.ID,
			&device.UserID,
			&device.OS,
			&device.Browser,
			&device.DeviceType,
			&device.IPAddress,
			&device.UserAgent,
			&device.LastUsedAt,
			&device.Region,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device row: %w", err)
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over device rows: %w", err)
	}

	return devices, nil
}

func (r *DeviceRepository) DeleteByIPAddress(ctx context.Context, ipAddress string) error {
	query := `DELETE FROM auth.devices WHERE ip_address = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to delete devices by IP address: %w", err)
	}

	return nil
}


package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DeviceService interface {
	CreateDevice(device *models.Device) (uuid.UUID, error)
	GetDevice(id string) (models.Device, error)
	UpdateDevice(device *models.Device) (models.Device, error)
	DeleteDevice(id string) error
	GetDevicesByUser(id string) ([]models.Device, error)
	DeleteDevicesByUser(id string) error
	UpdateDeviceLastUsedAt(id string, lastUsedAt time.Time) error
	CountDevicesByUser(id string) (int, error)
	DeleteExpiredDevices(expirationTime time.Time) error
	GetDevicesByIPAddress(ipAddress string) ([]models.Device, error)
	DeleteDevicesByIPAddress(ipAddress string) error
}

type deviceService struct {
	deviceRepo repository.DeviceRepository
}

func NewDeviceService() DeviceService {
	return &deviceService{
		deviceRepo: *repository.NewDeviceRepository(),
	}
}

func (s *deviceService) CreateDevice(device *models.Device) (uuid.UUID, error) {
	ctx := context.Background()
	createdDeviceID, err := s.deviceRepo.Insert(ctx, device)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create device: %w", err)
	}
	return createdDeviceID, nil
}

func (s *deviceService) GetDevice(id string) (models.Device, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Device{}, fmt.Errorf("invalid device ID: %w", err)
	}
	device, err := s.deviceRepo.Query(ctx, parsedID)
	if err != nil {
		return models.Device{}, fmt.Errorf("failed to get device: %w", err)
	}
	return device, nil
}

func (s *deviceService) UpdateDevice(device *models.Device) (models.Device, error) {
	ctx := context.Background()
	updatedDevice, err := s.deviceRepo.Update(ctx, device)
	if err != nil {
		return models.Device{}, fmt.Errorf("failed to update device: %w", err)
	}
	return updatedDevice, nil
}

func (s *deviceService) DeleteDevice(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid device ID: %w", err)
	}
	err = s.deviceRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}
	return nil
}

func (s *deviceService) GetDevicesByUser(id string) ([]models.Device, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	devices, err := s.deviceRepo.QueryByUser(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by user ID: %w", err)
	}
	return devices, nil
}

func (s *deviceService) DeleteDevicesByUser(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	err = s.deviceRepo.DeleteByUser(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete devices by user ID: %w", err)
	}
	return nil
}

func (s *deviceService) UpdateDeviceLastUsedAt(id string, lastUsedAt time.Time) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid device ID: %w", err)
	}
	err = s.deviceRepo.UpdateLastUsedAt(ctx, parsedID, lastUsedAt)
	if err != nil {
		return fmt.Errorf("failed to update device last used at: %w", err)
	}
	return nil}

func (s *deviceService) CountDevicesByUser(id string) (int, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}
	count, err := s.deviceRepo.CountDevicesByUser(ctx, parsedID)
	if err != nil {
		return 0, fmt.Errorf("failed to count devices by user ID: %w", err)
	}
	return count, nil
}

func (s *deviceService) DeleteExpiredDevices(expirationTime time.Time) error {
	ctx := context.Background()
	err := s.deviceRepo.DeleteExpiredDevices(ctx, expirationTime)
	if err != nil {
		return fmt.Errorf("failed to delete expired devices: %w", err)
	}
	return nil
}

func (s *deviceService) GetDevicesByIPAddress(ipAddress string) ([]models.Device, error) {
	ctx := context.Background()
	devices, err := s.deviceRepo.QueryByIPAddress(ctx, ipAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by IP address: %w", err)
	}
	return devices, nil
}

func (s *deviceService) DeleteDevicesByIPAddress(ipAddress string) error {
	ctx := context.Background()
	err := s.deviceRepo.DeleteByIPAddress(ctx, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to delete devices by IP address: %w", err)
	}
	return nil
}

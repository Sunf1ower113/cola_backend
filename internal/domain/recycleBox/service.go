package recycleBox

import (
	"context"
)

type ServiceRecycleBox interface {
	GetRecycleBox(ctx context.Context, id int64) (*RecycleBox, error)
	CreateRecycleBox(ctx context.Context, dto *CreateRecycleBoxDTO) (*RecycleBox, error)
	UpdateRecycleBox(ctx context.Context, id int64, dto *UpdateRecycleBoxDTO) (*RecycleBox, error)
	AddBottle(ctx context.Context, boxId int64) (*RecycleBox, error)
	AddBottleWithPoints(ctx context.Context, boxId int64, userId int64) (*RecycleBox, error)
}

type serviceRecycleBox struct {
	storage RecycleBoxStorage
}

func NewRecycleBoxService(storage RecycleBoxStorage) ServiceRecycleBox {
	return &serviceRecycleBox{
		storage: storage,
	}
}

// GetRecycleBox retrieves a recycle box by ID
func (s *serviceRecycleBox) GetRecycleBox(ctx context.Context, id int64) (*RecycleBox, error) {
	return s.storage.GetRecycleBox(id)
}

// CreateRecycleBox creates a new recycle box
func (s *serviceRecycleBox) CreateRecycleBox(ctx context.Context, dto *CreateRecycleBoxDTO) (*RecycleBox, error) {
	return s.storage.CreateRecycleBox(dto)
}

// UpdateRecycleBox updates an existing recycle box's details
func (s *serviceRecycleBox) UpdateRecycleBox(ctx context.Context, id int64, dto *UpdateRecycleBoxDTO) (*RecycleBox, error) {
	return s.storage.UpdateRecycleBox(id, dto)
}

// AddBottle increments bottle count in the recycle box without awarding points
func (s *serviceRecycleBox) AddBottle(ctx context.Context, boxId int64) (*RecycleBox, error) {
	return s.storage.AddBottle(boxId)
}

// AddBottleWithPoints increments bottle count in the recycle box and awards points to the user
func (s *serviceRecycleBox) AddBottleWithPoints(ctx context.Context, boxId int64, userId int64) (*RecycleBox, error) {
	// Attempt to add a bottle to the recycle box
	rb, err := s.storage.AddBottleWithPoints(boxId, userId)
	if err != nil {
		return nil, err
	}

	return rb, nil
}

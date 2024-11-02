package recycleBox

import (
	"auth-api/internal/domain/recycleBox"
	customError "auth-api/internal/error"
	"database/sql"
	"errors"
	"log"
)

const (
	points = 100
)

func NewRecycleBoxStorage(db *sql.DB) recycleBox.RecycleBoxStorage {
	return &storageRecycleBox{
		db: db,
	}
}

type storageRecycleBox struct {
	db *sql.DB
}

func (s *storageRecycleBox) GetRecycleBox(id int64) (*recycleBox.RecycleBox, error) {
	rb := &recycleBox.RecycleBox{}
	q := `SELECT id, title, address, capacity, count FROM recycle_boxes WHERE id = ?`
	row := s.db.QueryRow(q, id)
	if err := row.Scan(&rb.Id, &rb.Title, &rb.Address, &rb.Capacity, &rb.Count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.NotFoundError
		}
		return nil, err
	}
	return rb, nil
}

// CreateRecycleBox inserts a new RecycleBox using a DTO
func (s *storageRecycleBox) CreateRecycleBox(dto *recycleBox.CreateRecycleBoxDTO) (*recycleBox.RecycleBox, error) {
	q := `INSERT INTO recycle_boxes(title, address, capacity, count) VALUES (?, ?, ?, 0)`
	result, err := s.db.Exec(q, dto.Title, dto.Address, dto.Capacity)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &recycleBox.RecycleBox{
		Id:       id,
		Title:    dto.Title,
		Address:  dto.Address,
		Capacity: dto.Capacity,
		Count:    0,
	}, nil
}

// UpdateRecycleBox updates an existing RecycleBox based on the provided DTO
func (s *storageRecycleBox) UpdateRecycleBox(id int64, dto *recycleBox.UpdateRecycleBoxDTO) (*recycleBox.RecycleBox, error) {
	q := `UPDATE recycle_boxes SET title = ?, address = ?, capacity = ?, count = ? WHERE id = ?`
	_, err := s.db.Exec(q, dto.Title, dto.Address, dto.Capacity, dto.Count, id)
	if err != nil {
		return nil, err
	}

	return s.GetRecycleBox(id)
}

func (s *storageRecycleBox) FlushRecycleBox(id int64) (*recycleBox.RecycleBox, error) {
	q := `UPDATE recycle_boxes SET count = 0 WHERE id = ?`
	_, err := s.db.Exec(q, id)
	if err != nil {
		return nil, err
	}

	return s.GetRecycleBox(id)
}

func (s *storageRecycleBox) AddBottleWithPoints(boxId int64, userId int64) (*recycleBox.RecycleBox, error) {
	// First, call AddBottle to increment the count in the recycle box
	rb, err := s.AddBottle(boxId)
	if err != nil {
		return nil, err
	}
	log.Println("userid: ", userId)
	// Award points to the user
	q := `UPDATE users SET points = points + ? WHERE user_id = ?`
	_, err = s.db.Exec(q, 100, userId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return rb, nil
}

func (s *storageRecycleBox) AddBottle(id int64) (*recycleBox.RecycleBox, error) {
	// First, retrieve the current count and capacity to check if the box is full
	rb := &recycleBox.RecycleBox{}
	qSelect := `SELECT id, title, address, capacity, count FROM recycle_boxes WHERE id = ?`
	row := s.db.QueryRow(qSelect, id)
	if err := row.Scan(&rb.Id, &rb.Title, &rb.Address, &rb.Capacity, &rb.Count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.NotFoundError
		}
		return nil, err
	}

	// Check if the RecycleBox is already at full capacity
	if rb.Count >= rb.Capacity {
		return nil, errors.New("recycle box is full")
	}

	// Increment the count
	qUpdate := `UPDATE recycle_boxes SET count = count + 1 WHERE id = ?`
	_, err := s.db.Exec(qUpdate, id)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated record and return it
	return s.GetRecycleBox(id)
}

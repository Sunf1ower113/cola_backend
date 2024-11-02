package recycleBox

type RecycleBoxStorage interface {
	GetRecycleBox(int64) (*RecycleBox, error)
	CreateRecycleBox(*CreateRecycleBoxDTO) (*RecycleBox, error)
	UpdateRecycleBox(int64, *UpdateRecycleBoxDTO) (*RecycleBox, error)
	FlushRecycleBox(int64) (*RecycleBox, error)
	AddBottle(int64) (*RecycleBox, error)
	AddBottleWithPoints(int64, int64) (*RecycleBox, error)
}

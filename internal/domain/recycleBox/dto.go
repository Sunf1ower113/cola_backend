package recycleBox

type CreateRecycleBoxDTO struct {
	Title    string `json:"title"`
	Address  string `json:"address"`
	Capacity int64  `json:"capacity"`
}
type UpdateRecycleBoxDTO struct {
	Title    string `json:"title"`
	Address  string `json:"address"`
	Capacity int64  `json:"capacity"`
	Count    int64  `json:"count"`
}

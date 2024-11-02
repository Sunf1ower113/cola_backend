package recycleBox

type RecycleBox struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Address  string `json:"address"`
	Capacity int64  `json:"capacity"`
	Count    int64  `json:"count"`
}

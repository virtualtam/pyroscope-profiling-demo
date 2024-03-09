package restaurant

type Restaurant struct {
	Model

	Name string `json:"name"`
	Menu Menu   `json:"menu"`
}

package vo

type SetMealDishesUserVO struct {
	Copies      int32  `json:"copies"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Name        string `json:"name"`
}

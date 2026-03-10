package models

type TaskDto struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
}

func (t *Task) ToDto() *TaskDto {
	return &TaskDto{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

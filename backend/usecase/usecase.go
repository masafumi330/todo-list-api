package usecase

import "sync"

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TodoUsecase struct {
	mu     sync.Mutex
	nextID int
	todos  map[int]Todo
}

func NewTodoUsecase() *TodoUsecase {
	return &TodoUsecase{
		nextID: 1,
		todos:  make(map[int]Todo),
	}
}

func (uc *TodoUsecase) Create(title, description string) Todo {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	todo := Todo{
		ID:          uc.nextID,
		Title:       title,
		Description: description,
	}

	uc.todos[todo.ID] = todo
	uc.nextID++

	return todo
}

func (uc *TodoUsecase) Update(id int, title, description string) (Todo, bool) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	existing, ok := uc.todos[id]
	if !ok {
		return Todo{}, false
	}

	existing.Title = title
	existing.Description = description
	uc.todos[id] = existing

	return existing, true
}

func (uc *TodoUsecase) Delete(id int) bool {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.todos[id]; !ok {
		return false
	}

	delete(uc.todos, id)
	return true
}

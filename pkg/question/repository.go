package question

import "context"

// Respository interface, use it to define and expose all the methods available for our type Question
type Respository interface {
	GetAll(ctx context.Context) ([]Question, error)
	Create(ctx context.Context, company *Question) error
	Update(ctx context.Context, id uint, question Question) error
	Delete(ctx context.Context, id uint) error
}

package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jrmanes/go-toggl/pkg/question"
)

// QuestionRepository will be a bridge to data which will give us access to DBs
type QuestionRepository struct {
	Data *Data
}

// GetAll implement a user repository against infrastructure
func (qu *QuestionRepository) GetAll(ctx context.Context) ([]question.Question, error) {
	q := `
	SELECT id, body, created_at, updated_at
	FROM questions;
	`
	//SELECT id, body, options
	//FROM questions;
	//
	rows, err := qu.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []question.Question
	for rows.Next() {
		var qu question.Question
		rows.Scan(&qu.ID, &qu.Body, &qu.Options)
		questions = append(questions, qu)
	}

	return questions, nil
}


// Create adds a new user.
func (qu *QuestionRepository) Create(ctx context.Context, q *question.Question) error {
	query := `
	INSERT INTO questions (body, created_at, updated_at)
	VALUES ($1, $2, $3);
	`

	stmt, err := qu.Data.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	stmt.QueryRowContext(ctx, q.Body, time.Now(), time.Now()).Scan(q.ID)

	lastid := qu.getBookLastID()

	//create a loop in order to add the tags into the tags table
	fmt.Println("optiones:", q.Options)
	for i := 0; i < len(q.Options); i++ {
		query2 := `
		INSERT INTO options (id, body, correct)
		VALUES ((select id from questions where id=$1), $2, $3);
		`

		//prepare the statement
		stmt2, err := qu.Data.DB.PrepareContext(ctx, query2)
		if err != nil {
			return err
		}

		defer stmt2.Close()

		stmt2.QueryRowContext(ctx, lastid, q.Options[i].Body, q.Options[i].Correct).Scan(&q.ID)
	}

	query3 := `
	SELECT id, body, options( 
		select body, correct from options where id = $1) 
	FROM questions 
	WHERE id=$1;
	`
	row := stmt.QueryRow(query3, lastid)
	err = row.Scan(&q.ID)
	if err != nil {
		return  err
	}

	return nil
}

func (qu *QuestionRepository) getBookLastID() int {
	var id int

	query := `
	SELECT IFNULL(MAX(id), 0) AS id
	FROM questions`

//	defer qu.Data.DB.Close()

	err := qu.Data.DB.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	return id
}

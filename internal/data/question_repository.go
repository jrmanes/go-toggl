package data

import (
	"context"
	"fmt"

	"github.com/jrmanes/go-toggl/pkg/question"
	"github.com/lib/pq"
)

// QuestionRepository will be a bridge to data which will give us access to DBs
type QuestionRepository struct {
	Data *Data
}

// GetAll implement a question repository against infrastructure
func (qu *QuestionRepository) GetAll(ctx context.Context) ([]question.Question, error) {
	q := `
	SELECT Q.id, 
	       Q.body
		FROM questions AS Q
		GROUP BY Q.id;
	`

	rows, err := qu.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []question.Question

	// reorder tags on array
	for rows.Next() {
		var q question.Question
		//var tagsArray []string
		rows.Scan(&q.ID, &q.Body, pq.Array(&q.Options))
		fmt.Println(q.ID)
		// split values into an array
		//for i := range q.Options {
		//	tagsArray = strings.Split(q.Options[i].Body, ", ")
		//}
		//
		//// convert original array into the new separated one
		//q.Options = tagsArray

		// in the following line, we provide the time since the offer was created

		questions = append(questions, q)
	}
	return questions, nil
}


// Create adds a new question.
func (qu *QuestionRepository) Create(ctx context.Context, q *question.Question) error {
	query := `
	INSERT INTO questions (body)
	VALUES ($1)
	RETURNING id;
	`

	stmt, err := qu.Data.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, q.Body)
	err = row.Scan(&q.ID)
	if err != nil {
		return err
	}

	lastid := qu.getLastID()

	//create a loop in order to add the tags into the tags table
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

	return nil
}

func (qu *QuestionRepository) getLastID() int {
	var id int

	query := `
	SELECT MAX(id) AS id
	FROM questions`

	err := qu.Data.DB.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

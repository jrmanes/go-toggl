package data

import (
	"context"

	"github.com/jrmanes/go-toggl/pkg/question"
	_ "github.com/lib/pq"
)

// QuestionRepository will be a bridge to data which will give us access to DBs
type QuestionRepository struct {
	Data *Data
}

// GetAll implement a question repository against infrastructure
func (qu *QuestionRepository) GetAll(ctx context.Context) ([]question.Question, error) {
	q := `
		SELECT DISTINCT Q.id,
		       Q.body,
		        array(
		           SELECT DISTINCT ROW (body, correct)
		           FROM "options"
		         ) as "options"
		FROM questions AS Q, "options" AS O
		GROUP BY Q.id, O.id;
	`

	rows, err := qu.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []question.Question

	for rows.Next() {
		//var q question.Question
		//var tagsArray question.
		//rows.Scan(&q.ID, &q.Body, pq.Array(&q.Options))
		//fmt.Println(q)
		//// split values into an array
		//for i := range q.Options {
		//	tagsArray = strings.Split(q.Options[i].Body, ", ")
		//}
		//
		//// convert original array into the new separated one
		//q.Options.Body = tagsArray
		//
		//// in the following line, we provide the time since the question was created
		//
		//questions = append(questions, q)
	}
	return questions, nil
}

// Create adds a new question.
func (qu *QuestionRepository) Create(ctx context.Context, q *question.Question) error {
	queryCreate := `
	INSERT INTO questions (body)
	VALUES ($1)
	RETURNING id;
	`

	stmt, err := qu.Data.DB.PrepareContext(ctx, queryCreate)
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

// Update updates a question by id.
func (qu *QuestionRepository) Update(ctx context.Context, id uint, q question.Question) error {
	queryUpdate := `
	UPDATE questions set body=$1
	WHERE id=$2
	RETURNING id;
	`

	stmt, err := qu.Data.DB.PrepareContext(ctx, queryUpdate)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, q.Body, id)
	if err != nil {
		return err
	}

	// Clean up all the options
	qu.DeleteFromOptions(ctx, id)

	//create a loop in order to add the tags into the tags table
	for i := 0; i < len(q.Options); i++ {
		queryInsertIntoOptions := `
		INSERT INTO options (id, body, correct)
		VALUES ((select id from questions where id=$1), $2, $3);
		`
		stmt, err := qu.Data.DB.PrepareContext(ctx, queryInsertIntoOptions)
		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.ExecContext(ctx, id, q.Options[i].Body, q.Options[i].Correct)
		if err != nil {
			return err
		}
	}


	return nil
}

// Delete removes a question by id.
func (qu *QuestionRepository) Delete(ctx context.Context, id uint) error {
	queryDelete := `DELETE FROM questions WHERE id=$1;`

	stmt, err := qu.Data.DB.PrepareContext(ctx, queryDelete)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromOptions removes a question by id.
func (qu *QuestionRepository) DeleteFromOptions(ctx context.Context, id uint) error {
	queryDeleteFromOptions := `DELETE FROM "options" WHERE id=$1;`

	stmt, err := qu.Data.DB.PrepareContext(ctx, queryDeleteFromOptions)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// getLastID assert function to get the last id inserted
func (qu *QuestionRepository) getLastID() int {
	var id int

	query := `
	SELECT MAX(id) AS id
	FROM questions
	`

	err := qu.Data.DB.QueryRow(query).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

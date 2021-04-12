package data

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jrmanes/go-toggl/pkg/question"
	_ "github.com/lib/pq"
)

// QuestionRepository will be a bridge to data which will give us access to DBs
type QuestionRepository struct {
	Data *Data
}

// TODO: Pending to generate the array of options to return it when get all data
// GetAll implement a question repository against infrastructure
func (qu *QuestionRepository) GetAll(ctx context.Context) ([]question.Question, error) {
	q := `
	SELECT DISTINCT Q.id,
		   Q.body,
		   (
			   SELECT ARRAY (
				   SELECT  body
				   FROM options
				   WHERE options.id = Q.id
			  )
		   ) AS "options_body",
		   (
			   SELECT ARRAY (
				   SELECT  correct
				   FROM options
				   WHERE options.id = Q.id
			  )
		   ) AS "options_correct"
	FROM questions AS Q
			 INNER JOIN "options"
						ON Q.id = options.id
	GROUP BY Q.id, options.body, Q.body, Q.id, options.correct
	ORDER BY Q.id;
	`

	rows, err := qu.Data.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var questions []question.Question
	for rows.Next() {
		var qu question.Question
		var op question.Options

		err := rows.Scan(&qu.ID, &qu.Body, &op.Body, &op.Correct)
		if err != nil {
			log.Fatal(err)
		}
		var body string = op.Body
		for _, field := range split(body, '{') {
			for _, field2 := range split(field, ',') {
				for _, field3 := range split(field2, '}') {
					fmt.Println(field3)
					//qu.Options[0].Body = field3
					//questions = append(questions, qu)
				}
			}
		}
		for i := 0; i < (len(op.Correct)); i++ {
			for _, field := range split(string(op.Correct[i]), '{') {
				for _, field2 := range split(field, ',') {
					for _, field3 := range split(field2, '}') {
						if field3 == "t" {
							field3 = "true"
						}
						if field3 == "" {
							field3 = "false"
						}
						correct, err := strconv.ParseBool(field3)
						if err != nil {
							log.Fatal("error parse bool", err)
						}
						fmt.Println(correct)

						//qu.Options[0].Body = append(qu.Options, &op.{body: field3})
						//qu.Options[0].Body = field3
						//questions = append(questions, qu)

					}
				}
			}
		}

		questions = append(questions, qu)
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

// split return divided strings
func split(tosplit string, sep rune) []string {
	var fields []string

	last := 0
	for i, c := range tosplit {
		if c == sep {
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1
		}
	}

	// last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}

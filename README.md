# go-toggl


### API Requirements

Questions have a simple structure.
Each question has a body. Then there are two or more options.
Each option has a body as well and a boolean attribute that defines whether the option is correct. 
At least one of the options is correct.
Below is a JSON representation of a sample question.

```json
{
  "body": "Where does the sun set?",
  "options": [
    {
      "body": "East",
      "correct": false
    },
    {
      "body": "West",
      "correct": true
    }
  ]
}
```

### Listing questions

Endpoint:
```http
http://localhost:3000/q
```
METHOD:
```http
GET
```

The first endpoint should return a list of all questions in the database.
The order of questions and options inside questions should be stable, i.e. not change on every request. The whole question, including the options is returned.

For example, the response could look like this:

```json5
[
  {
    "body": "Where does the sun set?",
    "options": [
      {
        "body": "East",
        "correct": false
      },
      // other options...
    ]
  },
  {
    "body": "What is the answer to the ultimate question of life, the universe, and everything?",
    // rest of the question...
  },
  {
    "body": "But what is the ultimate question?",
    // rest of the question...
  }
]
```

### Creating a new question
Endpoint:
```http
http://localhost:3000/q
```
METHOD:
```http
POST
```

The second endpoint creates a new question in the database and then returns it in the response. The request body contains the question in JSON. The order of options in the request body should be stored as well and the same order should be returned by the API from all requests.

For example, for the request containing the following JSON, the server would return the question show above, in the Questions section.

```json
{
  "body": "Where does the sun set?",
  "options": [
    {
      "body": "East",
      "correct": false
    },
    {
      "body": "West",
      "correct": true
    }
  ]
}
```

### Updating a question
Endpoint:
```http
http://localhost:3000/q/{id}
```
METHOD:
```http
PUT
```

The third endpoint updates an existing question and returns the updated question in the response.
The whole question is included in the request body, including all attributes. 
The question to be updated should be identified in the request URL. 
The order of options in the request body should be stored the same way as explained in the create endpoint above.

For example, to change the question from before to ask about sunrise, we would send the following JSON.

```json
{
  "body": "Where does the sun rise?",
  "options": [
    {
      "body": "East",
      "correct": true
    },
    {
      "body": "West",
      "correct": false
    }
  ]
}
```

## Basic requirements

Your solution should meet all these requirements.

- [ ] Endpoint that returns a list of all questions
- [X] Endpoint that allows to add a new question
- [X] Endpoint that allows to update an existing question
- [X] Question data is stored in a SQLite database with a **normalised** schema
  Use **PostgreSQL** with [migrations](https://github.com/jrmanes/go-toggl/tree/main/internal/data/db/migrations)
  Running in a [docker-compose](https://github.com/jrmanes/go-toggl/blob/main/infra/docker/docker-compose.yml)
- [X] The order of questions and options is stable, not random
- [X] The `PORT` environment variable is used as the port number for the server, defaulting to 3000

## Bonus requirements

These requirements are not required, but feel free to complete some of them if they seem interesting, or to come up with your own :)

- [X] Endpoint that allows to delete existing questions
  Endpoint:
```http
http://localhost:3000/q/{id}
```
METHOD:
```http
DELETE
```
- [ ] Pagination for the list endpoint

  This can be in the form of basic offset pagination, or seek pagination. The difference is explained in [this post](https://web.archive.org/web/20210205081113/https://taylorbrazelton.com/posts/2019/03/offset-vs-seek-pagination/).

- [X] JWT authentication mechanism

  Clients are required to send a JSON Web Token that identifies the user in some way. The API returns only questions that belong to the authenticated user. Endpoint for generating tokens is not needed, we can generate them through [jwt.io](https://jwt.io/).

- [ ] Use GraphQL instead of REST to implement the API

  Define a schema for the API that covers the basic requirements and implement all queries and resolvers. You do not need to implement the REST API if you choose to do this.


## Setup the project


In order to make it easy, there is a Makefile with different acctions, to start the project, just execute:
```Makefile
make setup_project
```

### Apply migrations to database
In order to create the **first schema** of the database, just execute:
```Makefile
make migrate_up
```
That command will execute the migrations against the database.



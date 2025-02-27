# Backend Server Technical Specs

## Business Goal

A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps


## Technical Requirements

- Backend to be built with Golang
- Use SQLite3 as the database
- Use Gin or Fiber as the web framework -> Decided on Gin (see below)
- We should use a task runner (TBD which one)
- The API format will be JSON
- There will be no authentication/authorization required
- We will assume a single user for the purposes of this project




### Deciding on Gin or Fiber

After a bit of research and discussion with Anthropic's Claude, Gin and Fiber
are popular options. I decided to go with Gin because it has a more Go-idiomatic
approach. Although it may be slower than Fiber, I believe that this is the
better choice for this project because it aligns better with my desire to learn
and expand my GO knowledge and its idiosyncrasies.

Summary:

```
Performance:

- Fiber is built on Fasthttp and typically benchmarks faster than Gin
- Gin uses the standard net/http package which is more battle-tested but slightly slower

Syntax and Learning Curve:

- Fiber has Express-like syntax that might feel more familiar if you're coming from Node.js
- Gin has a more Go-idiomatic approach that aligns better with Go's philosophy

Middleware Ecosystem:

- Gin has a larger community and more third-party middleware available
- Fiber has a growing ecosystem but currently has fewer options

Stability:

- Gin is more mature and stable with fewer breaking changes
- Fiber is newer and still evolving, which can mean occasional breaking changes

When to choose Fiber:

- You prefer Express/FastAPI-like syntax
- Absolute maximum performance is a priority
- You're building a new project and don't mind some API changes down the road

When to choose Gin:

- You want a more established, stable framework
- You prefer Go-idiomatic code patterns
- You need access to a wide range of third-party middleware
- You're working on production-critical services where library maturity matters
```

## Database schema

We have the following tables:

- words - Stores individual Portuguese vocabulary words
    - `id` (Primary Key): Unique identifier for each word
    - `portuguese` (String, Required): The word written in Portuguese
    - `english` (String, Required): English translation of the word

- groups - thematic word groups
    - `id` (Primary Key): Unique identifier for each group
    - `name` (String, Required): Name of the group

- words_groups - join-table enabling many-to-many relationship between words and groups
    - `id` (Primary Key): Unique identifier for each association
    - `word_id` (Foreign Key): References words.id
    - `group_id` (Foreign Key): References groups.id

- study_activities - a specific study activity, linking a study session to a
group
    - `id` (Primary Key): Unique identifier for each activity
    - `study_session_id` (Foreign Key): References study_sessions.id
    - `group_id` (Foreign Key): References groups.id
    - `created_at` (Timestamp, Default: Current Time): When the activity was created

- study_sessions - records of study sessions grouping word_review_items
    - `id` (Primary Key): Unique identifier for each session
    - `group_id` (Foreign Key): References groups.id
    - `study_activity_id` (Foreign Key): References study_activities.id
    - `created_at` (Timestamp, Default: Current Time): When the session was created

- word_review_items - a record of word practice, determining if the word was
correct or not
    - `id` (Primary Key): Unique identifier for each review item
    - `word_id` (Foreign Key): References words.id
    - `study_session_id` (Foreign Key): References study_sessions.id
    - `correct` (Boolean): Whether the word was practiced correctly
    - `created_at` (Timestamp, Default: Current Time): When the review item was created

## Endpoints

### GET /api/dashboard/last_study_session
Returns information about the most recent study session.

#### JSON Response
```json
{
  "id": 123,
  "group_id": 456,
  "created_at": "2025-02-08T17:20:23-05:00",
  "study_activity_id": 789,
  "group_id": 456,
  "group_name": "Basic Greetings"
}
```

### GET /api/dashboard/study_progress
Returns study progress statistics.

Note: that the frontend will determine progress bar based on total words studied and total available words.

#### JSON Response

```json
{
  "total_words_studied": 3,
  "total_available_words": 124,
}
```

### GET /api/dashboard/quick-stats

Returns quick overview statistics.

#### JSON Response
```json
{
  "success_rate": 80.0,
  "total_study_sessions": 4,
  "total_active_groups": 3,
  "study_streak_days": 4
}
```

### GET /api/study_activities/:id

#### JSON Response
```json
{
  "id": 1,
  "name": "Vocabulary Quiz",
  "thumbnail_url": "https://example.com/thumbnail.jpg",
  "description": "Practice your vocabulary with flashcards"
}
```

### GET /api/study_activities/:id/study_sessions

- pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 20
  }
}
```

### POST /api/study_activities

#### Request Params
- group_id integer
- study_activity_id integer

#### JSON Response
{
  "id": 124,
  "group_id": 123
}

### GET /api/words

- pagination with 100 items per page

#### JSON Response
```json
{
  "items": [
    {
      "portuguese": "olá",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 500,
    "items_per_page": 100
  }
}
```

### GET /api/words/:id

#### JSON Response
```json
{
  "portuguese": "olá",
  "english": "hello",
  "stats": {
    "correct_count": 5,
    "wrong_count": 2
  },
  "groups": [
    {
      "id": 1,
      "name": "Basic Greetings"
    }
  ]
}
```

### GET /api/groups
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "id": 1,
      "name": "Basic Greetings",
      "word_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 10,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id
#### JSON Response
```json
{
  "id": 1,
  "name": "Basic Greetings",
  "stats": {
    "total_word_count": 20
  }
}
```

### GET /api/groups/:id/words
#### JSON Response
```json
{
  "items": [
    {
      "portuguese": "olá",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 100
  }
}
```

### GET /api/groups/:id/study_sessions
#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 5,
    "items_per_page": 100
  }
}
```

### GET /api/study_sessions
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "id": 123,
      "activity_name": "Vocabulary Quiz",
      "group_name": "Basic Greetings",
      "start_time": "2025-02-08T17:20:23-05:00",
      "end_time": "2025-02-08T17:30:23-05:00",
      "review_items_count": 20
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 100,
    "items_per_page": 100
  }
}
```

### GET /api/study_sessions/:id
#### JSON Response
```json
{
  "id": 123,
  "activity_name": "Vocabulary Quiz",
  "group_name": "Basic Greetings",
  "start_time": "2025-02-08T17:20:23-05:00",
  "end_time": "2025-02-08T17:30:23-05:00",
  "review_items_count": 20
}
```

### GET /api/study_sessions/:id/words
- pagination with 100 items per page
#### JSON Response
```json
{
  "items": [
    {
      "japanese": "こんにちは",
      "romaji": "konnichiwa",
      "english": "hello",
      "correct_count": 5,
      "wrong_count": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_items": 20,
    "items_per_page": 100
  }
}
```

### POST /api/reset_history
#### JSON Response
```json
{
  "success": true,
  "message": "Study history has been reset"
}
```

### POST /api/full_reset
#### JSON Response
```json
{
  "success": true,
  "message": "System has been fully reset"
}
```

### POST /api/study_sessions/:id/words/:word_id/review
#### Request Params
- id (study_session_id) integer
- word_id integer
- correct boolean

#### Request Payload
```json
{
  "correct": true
}
```

#### JSON Response
```json
{
  "success": true,
  "word_id": 1,
  "study_session_id": 123,
  "correct": true,
  "created_at": "2025-02-08T17:33:07-05:00"
}
```

## Required Automated Tasks

Lets list out possible tasks we need for our lang portal.

### Initialize Database
This task will initialize the sqlite database called `words.db

### Migrate Database
This task will run a series of migrations sql files on the database

Migrations live in the `migrations` folder.
The migration files will be run in order of their file name.
The file names should looks like this:

```sql
0001_init.sql
0002_create_words_table.sql
```

### Seed Data
This task will import json files and transform them into target data for our database.

All seed files live in the `seeds` folder.

In our task we should have DSL to specific each seed file and its expected group word name.

```json
[
  {
    "portuguese": "pagar",
    "english": "to pay",
  },
  ...
]
```
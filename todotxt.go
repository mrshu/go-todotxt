package todotxt

import (
        "time"
)

type Task struct {
        todo string
        priority string
        create_date time.Time
        contexts []string
        projects []string
}

type TaskList []Task


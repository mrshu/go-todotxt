package todotxt

import (
        "time"
        "fmt"
        "os"
        "bufio"
        "strings"
        "regexp"
)

type Task struct {
        todo string
        priority string
        create_date time.Time
        contexts []string
        projects []string
}

type TaskList []Task

func BuildTaskList (filename string) (TaskList) {

        var f, err = os.Open(filename)

        if err != nil {
                panic(err)
        }

        defer f.Close()

        var tasklist = TaskList{}

        scanner := bufio.NewScanner(f)

        for scanner.Scan() {
                var task = Task{}
                text := scanner.Text()
                splits := strings.Split(text, " ")
                date_regexp := "([\\d]{4})-([\\d]{2})-([\\d]{2})"

                if match, _ := regexp.MatchString(date_regexp, splits[0]); match {
                        t := fmt.Sprintf("%sT00:00:00Z", splits[0])
                        if date, e := time.Parse(time.RFC3339, t); e != nil {
                                panic(e)
                        } else {
                                task.create_date = date
                        }
                }

                fmt.Println(task.create_date)

                tasklist = append(tasklist, task)
        }

        if err := scanner.Err(); err != nil {
                panic(scanner.Err())
        }

        return tasklist
}


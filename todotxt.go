package todotxt

import (
        "time"
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

                match, _ := regexp.MatchString("[\\d]{4}-[\\d]{2}-[\\d]{2}", splits[0])
                if match {
                        task.create_date, _ = time.Parse(splits[0], "2013-12-30")
                }

                tasklist = append(tasklist, task)
        }

        if err := scanner.Err(); err != nil {
                panic(scanner.Err())
        }

        return tasklist
}


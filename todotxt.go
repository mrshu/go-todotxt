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
        priority byte
        create_date time.Time
        contexts []string
        projects []string
}

type TaskList []Task

func LoadTaskList (filename string) (TaskList) {

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

                head := splits[0]

                if (len(head) == 3) &&
                   (head[0] == '(') &&
                   (head[2] == ')') {
                        task.priority = head[1]
                        splits = splits[1:]
                }

                date_regexp := "([\\d]{4})-([\\d]{2})-([\\d]{2})"
                if match, _ := regexp.MatchString(date_regexp, splits[0]); match {
                        if date, e := time.Parse("2006-01-02", splits[0]); e != nil {
                                panic(e)
                        } else {
                                task.create_date = date
                        }

                        task.todo = strings.Join(splits[1:], " ")
                } else {
                        task.todo = strings.Join(splits[0:], " ")
                }

                context_regexp, _ := regexp.Compile("@[[:word:]]+")
                contexts := context_regexp.FindAllStringSubmatch(text, -1)
                if len(contexts) != 0 {
                        task.contexts = contexts[0]
                }

                project_regexp, _ := regexp.Compile("\\+[[:word:]]+")
                projects := project_regexp.FindAllStringSubmatch(text, -1)
                if len(projects) != 0 {
                        task.projects = projects[0]
                }

                tasklist = append(tasklist, task)
        }

        if err := scanner.Err(); err != nil {
                panic(scanner.Err())
        }

        return tasklist
}


func (tasks TaskList) Count() int {
        return len(tasks)
}

func (task Task) Text() string {
        return task.todo
}

func (task Task) Priority() byte {
        return task.priority
}

func (task Task) Contexts() []string {
        return task.contexts
}

func (task Task) Projects() []string {
        return task.projects
}

func (task Task) CreateDate() time.Time {
        return task.create_date
}

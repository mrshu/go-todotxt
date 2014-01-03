package todotxt

import (
        "time"
        "os"
        "bufio"
        "strings"
        "regexp"
        "sort"
)

type Task struct {
        todo string
        priority byte
        create_date time.Time
        contexts []string
        projects []string
        raw_todo string
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
                task.raw_todo = text

                splits := strings.Split(text, " ")

                head := splits[0]

                if (len(head) == 3) &&
                   (head[0] == '(') &&
                   (head[2] == ')') &&
                   (head[1] >= 65 && head[1] <= 90) { // checking if it's in range [A-Z]
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

func (tasks TaskList) Len() int {
        return len(tasks)
}

type ByPriority TaskList
func (tasks ByPriority) Len() int {
        return len(tasks)
}
func (tasks ByPriority) Swap(i, j int) {
        tasks[i], tasks[j] = tasks[j], tasks[i]
}
func (tasks ByPriority) Less(i, j int) bool {
        return tasks[i].Priority() < tasks[j].Priority()
}
func (tasks TaskList) Sort() {
        sort.Sort(ByPriority(tasks))
}

type ByCreateDate TaskList
func (tasks ByCreateDate) Len() int {
        return len(tasks)
}
func (tasks ByCreateDate) Swap(i, j int) {
        tasks[i], tasks[j] = tasks[j], tasks[i]
}
func (tasks ByCreateDate) Less(i, j int) bool {
        t1 := tasks[i].CreateDate().Unix()
        t2 := tasks[j].CreateDate().Unix()

        // if the dates equal, let's use priority
        if t1 == t2 {
                return tasks[i].Priority() < tasks[j].Priority()
        } else {
                return t1 > t2
        }
}
func (tasks TaskList) SortByCreateDate() {
        sort.Sort(ByCreateDate(tasks))
}


func (task Task) Text() string {
        return task.todo
}

func (task Task) RawText() string {
        return task.raw_todo
}

func (task Task) Priority() byte {
        // if priority is not from [A-Z], let it be 94 (^)
        if task.priority < 65 || task.priority > 90 {
                return 94 // you know, ^
        } else {
                return task.priority
        }
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

package todotxt

import (
        "time"
        "os"
        "bufio"
        "strings"
        "regexp"
        "sort"
        "unicode"
)

type Task struct {
        id int
        todo string
        priority byte
        create_date time.Time
        contexts []string
        projects []string
        raw_todo string
        finished bool
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
        id := 0

        for scanner.Scan() {
                var task = Task{}
                text := scanner.Text()
                task.id = id
                task.raw_todo = text

                splits := strings.Split(text, " ")

                if text[0] == 'x' &&
                   text[1] == ' ' &&
                   !unicode.IsSpace(rune(text[2])) {
                        task.finished = true
                        splits = splits[1:]
                }

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
                id += 1
        }

        if err := scanner.Err(); err != nil {
                panic(scanner.Err())
        }

        return tasklist
}

type By func(t1, t2 Task) bool

func (by By) Sort(tasks TaskList) {
        ts := &taskSorter{
                tasks: tasks,
                by:    by,
        }
        sort.Sort(ts)
}

type taskSorter struct {
        tasks TaskList
        by func(t1, t2 Task) bool
}

func (s *taskSorter) Len() int {
        return len(s.tasks)
}

func (s *taskSorter) Swap(i, j int) {
        s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

func (s *taskSorter) Less(i, j int) bool {
        return s.by(s.tasks[i], s.tasks[j])
}

func (tasks TaskList) Len() int {
        return len(tasks)
}

func prioCmp(t1, t2 Task) bool {
        return t1.Priority() < t2.Priority()
}

func prioRevCmp(t1, t2 Task) bool {
        return t1.Priority() > t2.Priority()
}

func dateCmp(t1, t2 Task) bool {
        tm1 := t1.CreateDate().Unix()
        tm2 := t2.CreateDate().Unix()

        // if the dates equal, let's use priority
        if tm1 == tm2 {
                return prioCmp(t1, t2)
        } else {
                return tm1 > tm2
        }
}

func dateRevCmp(t1, t2 Task) bool {
        tm1 := t1.CreateDate().Unix()
        tm2 := t2.CreateDate().Unix()

        // if the dates equal, let's use priority
        if tm1 == tm2 {
                return prioCmp(t1, t2)
        } else {
                return tm1 < tm2
        }
}

func lenCmp(t1, t2 Task) bool {
        tl1 := len(t1.raw_todo)
        tl2 := len(t2.raw_todo)
        if tl1 == tl2 {
                return prioCmp(t1, t2)
        } else {
                return tl1 < tl2
        }
}

func lenRevCmp(t1, t2 Task) bool {
        tl1 := len(t1.raw_todo)
        tl2 := len(t2.raw_todo)
        if tl1 == tl2 {
                return prioCmp(t1, t2)
        } else {
                return tl1 > tl2
        }
}

func idCmp(t1, t2 Task) bool {
        return t1.Id() < t2.Id()
}

func (tasks TaskList) Sort(by string) {
        switch by {
        default:
        case "prio":
                By(prioCmp).Sort(tasks)
        case "prio-rev":
                By(prioRevCmp).Sort(tasks)
        case "date":
                By(dateCmp).Sort(tasks)
        case "date-rev":
                By(dateRevCmp).Sort(tasks)
        case "len":
                By(lenCmp).Sort(tasks)
        case "len-rev":
                By(lenRevCmp).Sort(tasks)
        case "id":
                By(idCmp).Sort(tasks)
        }
}

func (tasks TaskList) Save(filename string) {
        tasks.Sort("id")

        f, err := os.Create(filename)
        if err != nil {
                panic(err)
        }

        defer f.Close()

        for _, task := range tasks {
                f.WriteString(task.RawText() + "\n")
        }
        f.Sync()
}

func (task Task) Id() int {
        return task.id
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

func (task Task) Finished() bool {
        return task.finished
}

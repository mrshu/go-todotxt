package todotxt

import "testing"

func TestLoadTaskList (t *testing.T) {
        tasklist := LoadTaskList("todo.txt")
        if tasklist.Count() != 7 {
                t.Errorf("Something went wrong with LoadTaskList: is %v, want %v, list %v",
                          tasklist.Count(), 7, tasklist)
        }
}

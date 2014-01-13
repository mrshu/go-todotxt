package todotxt

import "testing"

func TestLoadTaskList (t *testing.T) {
        tasklist := LoadTaskList("todo.txt")
        if tasklist.Len() != 8 {
                t.Errorf("Something went wrong with LoadTaskList: is %v, want %v\n%v",
                          tasklist.Len(), 8, tasklist)
        }
}

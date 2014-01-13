package todotxt

import "testing"

func TestLoadTaskList (t *testing.T) {
        tasklist := LoadTaskList("todo.txt")
        if tasklist.Len() != 8 {
                t.Errorf("Something went wrong with LoadTaskList: is %v, want %v\n%v",
                          tasklist.Len(), 8, tasklist)
        }
}

func TestLoadTaskListNonExistent (t *testing.T) {
        defer func(){
                if r:=recover(); r!=nil {
                        // recovered
                } else {
                        t.Errorf("Something went seriously wrong")
                }
        }()
        tasklist := LoadTaskList("nonexistent-file.txt")

        t.Errorf("Something is still wrong %v", tasklist)
}

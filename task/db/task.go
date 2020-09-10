package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

//Boltdb related global variables
var db *bolt.DB
var rootBucket = "DB"
var taskBucket = "TASK_TODO"
var completedTaskBucket = "COMPLETED_TASK"

//Struct to store incompleted task i.e task todo
type Task struct {
	ID   int
	Name string
}

//Struct to store completed task
type CompletedTask struct {
	ID             int
	Name           string
	CompletionTime time.Time
}

//Setup bolt db and create specified buckets, return error if any
func SetupDB(dbPath string) error {
	_db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("could not open db, %v", err)
	}
	db = _db

	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte(rootBucket))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte(taskBucket))
		if err != nil {
			return fmt.Errorf("could not create task_todo  bucket: %v", err)
		}

		_, err = root.CreateBucketIfNotExists([]byte(completedTaskBucket))
		if err != nil {
			return fmt.Errorf("could not create completed_task  bucket: %v", err)
		}

		return nil
	})

	return err
}

// Add a task to TASK_TODO bucket in boltdb
func AddTask(task Task) (int, error) {
	var taskId int

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, rootBucket, taskBucket)
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		task.ID = int(id)
		taskId = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(itob(task.ID), buf)
	})

	if err != nil {
		return -1, fmt.Errorf("count not add task to TODO list, %v", err)
	}
	return taskId, nil
}

// Fetch all incompleted tasks from bolt db
func FetchAllTask() ([]Task, error) {
	tasks := []Task{}

	err := db.View(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, rootBucket, taskBucket)
		if err != nil {
			return err
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			task := Task{}
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			tasks = append(tasks, task)
		}

		return nil
	})

	if err != nil {
		return []Task{}, err
	}
	return tasks, nil
}

//Remove a task with given id from bolt db
func RemoveTask(id int) error {
	err := db.Update(func(tx *bolt.Tx) error {

		b, err := getBucket(tx, rootBucket, taskBucket)
		if err != nil {
			return fmt.Errorf("failed to remove ,%v", err)
		}

		err = b.Delete(itob(id))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

// Fetch all completed tasks from bolt db
func FetchAllCompletedTask() ([]CompletedTask, error) {
	completedTasks := []CompletedTask{}

	err := db.View(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, rootBucket, completedTaskBucket)
		if err != nil {
			return err
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			completedTask := CompletedTask{}
			err := json.Unmarshal(v, &completedTask)
			if err != nil {
				return err
			}
			completedTasks = append(completedTasks, completedTask)
		}

		return nil
	})

	if err != nil {
		return []CompletedTask{}, err
	}
	return completedTasks, nil
}

//Add a task to COMPLETED_TASK bucket in boltdb
func AddCompletedTask(task CompletedTask) (int, error) {
	var completedTaskId = -1
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := getBucket(tx, rootBucket, completedTaskBucket)
		if err != nil {
			return err
		}

		id, err := b.NextSequence()

		if err != nil {
			return err
		}

		task.ID = int(id)
		completedTaskId = int(id)

		// Marshal user data into bytes.
		buf, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put(itob(task.ID), buf)
	})

	return completedTaskId, err
}

// get reference to the nested bucket
func getBucket(tx *bolt.Tx, rootBucket string, childBucket string) (*bolt.Bucket, error) {

	b := tx.Bucket([]byte(rootBucket))
	if b == nil {
		return nil, fmt.Errorf("%s bucket doesn't exists", rootBucket)
	}

	b = b.Bucket([]byte(childBucket))
	if b == nil {
		return nil, fmt.Errorf("%s bucket doesn't exists", childBucket)
	}
	return b, nil

}

//itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//itob returns an int from 8-byte big endian representation.
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

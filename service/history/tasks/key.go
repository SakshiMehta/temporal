package tasks

import (
	"fmt"
	"math"
	"time"
)

var (
	DefaultFireTime         = time.Unix(0, 0).UTC()
	defaultFireTimeUnixNano = DefaultFireTime.UnixNano()
)

var (
	MinimumKey = NewKey(DefaultFireTime, 0)
	MaximumKey = NewKey(time.Unix(0, math.MaxInt64), math.MaxInt64)
)

type (
	Key struct {
		// FireTime is the scheduled time of the task
		FireTime time.Time
		// TaskID is the ID of the task
		TaskID int64
	}

	Keys []Key
)

func NewImmediateKey(taskID int64) Key {
	return Key{
		FireTime: DefaultFireTime,
		TaskID:   taskID,
	}
}

func NewKey(fireTime time.Time, taskID int64) Key {
	return Key{
		FireTime: fireTime,
		TaskID:   taskID,
	}
}

func ValidateKey(key Key) error {
	if key.FireTime.UnixNano() < defaultFireTimeUnixNano {
		return fmt.Errorf("task key fire time must have unix nano value >= 0, got %v", key.FireTime.UnixNano())
	}

	if key.TaskID < 0 {
		return fmt.Errorf("task key ID must >= 0, got %v", key.TaskID)
	}

	return nil
}

func (left Key) CompareTo(right Key) int {
	if left.FireTime.Before(right.FireTime) {
		return -1
	} else if left.FireTime.After(right.FireTime) {
		return 1
	}

	if left.TaskID < right.TaskID {
		return -1
	} else if left.TaskID > right.TaskID {
		return 1
	}
	return 0
}

func (k Key) Prev() Key {
	if k.TaskID == 0 {
		if k.FireTime.UnixNano() == 0 {
			panic("Key encountered negative underflow")
		}
		return NewKey(k.FireTime.Add(-time.Nanosecond), math.MaxInt64)
	}
	return NewKey(k.FireTime, k.TaskID-1)
}

func (k Key) Next() Key {
	if k.TaskID == math.MaxInt64 {
		if k.FireTime.UnixNano() == math.MaxInt64 {
			panic("Key encountered positive overflow")
		}
		return NewKey(k.FireTime.Add(time.Nanosecond), 0)
	}
	return NewKey(k.FireTime, k.TaskID+1)
}

func (k Key) Sub(subtrahend Key) Key {
	borrow := int64(0)
	differenceTaskID := k.TaskID - subtrahend.TaskID
	if differenceTaskID < 0 {
		borrow = 1
		differenceTaskID += MaximumKey.TaskID
	}

	fireTime := k.FireTime.UnixNano() - borrow
	subtrahendFireTime := subtrahend.FireTime.UnixNano()
	if fireTime < subtrahendFireTime {
		panic(fmt.Sprintf("Task key Sub encountered underflow: self: %v, subtrahend: %v", k, subtrahend))
	}

	return NewKey(
		time.Unix(0, fireTime-subtrahendFireTime).UTC(),
		int64(differenceTaskID),
	)
}

func MinKey(this Key, that Key) Key {
	if this.CompareTo(that) < 0 {
		return this
	}
	return that
}

func MaxKey(this Key, that Key) Key {
	if this.CompareTo(that) < 0 {
		return that
	}
	return this
}

// Len implements sort.Interface
func (s Keys) Len() int {
	return len(s)
}

// Swap implements sort.Interface.
func (s Keys) Swap(
	this int,
	that int,
) {
	s[this], s[that] = s[that], s[this]
}

// Less implements sort.Interface
func (s Keys) Less(
	this int,
	that int,
) bool {
	return s[this].CompareTo(s[that]) < 0
}

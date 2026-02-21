package main

import "fmt"

type Message struct {
	Payload string
}

type Queue struct {
	messages []Message
}

func (q *Queue) Enqueue(msg Message) {
	q.messages = append(q.messages, msg)
}

func (q *Queue) Dequeue() (Message, error) {
	if len(q.messages) == 0 {
		return Message{}, fmt.Errorf("queue is empty")
	}
	msg := q.messages[len(q.messages)-1]
	q.messages = q.messages[:len(q.messages)-1]
	return msg, nil
}

func main() {
	var q = Queue{}
	fmt.Println(q)
	q.Enqueue(Message{"Hey1"})
	q.Enqueue(Message{"Hey2"})
	fmt.Println(q)
	val, err := q.Dequeue()
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	} else {
		fmt.Println(val.Payload)
	}

	val, err = q.Dequeue()
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
	} else {
		fmt.Println(val.Payload)
	}

	val, err = q.Dequeue()
	if err != nil {
		fmt.Printf("Error occured: %#v\n", err)
	} else {
		fmt.Println(val.Payload)
	}
}

package queue

import (
	"log"
	"time"

	"github.com/gocraft/work"
)

func (c *Context) SendEmail(job *work.Job) error {
	// Extract arguments:
	addr := job.ArgString("address")
	subject := job.ArgString("subject")
	if err := job.ArgError(); err != nil {
		return err
	}

	log.Println("Processing job " + job.ID + "...")

	time.Sleep(1 * time.Second)

	// Go ahead and send the email...
	// sendEmailTo(addr, subject)
	log.Println("address is " + addr)
	log.Println("subject is " + subject)

	return nil
}

func (c *Context) Export(job *work.Job) error {
	return nil
}

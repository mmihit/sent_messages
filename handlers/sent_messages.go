package handlers

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func (h *Handlers) SchedulerTask() error {
	nyc, err := time.LoadLocation("Africa/Casablanca")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithSeconds(), cron.WithLocation(nyc))

	c.AddFunc("*/10 * * * * *", h.SentMessages)

	fmt.Println("Cron running...")
	c.Start()
	return nil
}

func (h *Handlers) SentMessages() {
	nyc, err := time.LoadLocation("Africa/Casablanca")
	if err != nil {
		fmt.Println("error getting location time: ", err)
		return
	}

	day := time.Now().In(nyc).AddDate(0, 0, -1).Format("2006-01-02")

	patients, err := h.DB.GetPatientsFromScheduler(day)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, patient := range patients {
		fmt.Printf("sent message to client ID %d: %s %s, has number %s and %s\n", patient.ID, patient.FirstName, patient.LastName, patient.WhatsappNumber1, patient.WhatsappNumber2)
	}
}

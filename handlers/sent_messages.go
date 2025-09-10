package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"my_app/helper"
	"net/http"
	"strconv"
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
		err = sentMessageTo(patient, patient.WhatsappNumber1)
		if err == nil {
			//update message status in db
			fmt.Printf("\noriginal message sent to: %s %s, has this num %s\n", patient.FirstName, patient.LastName, patient.WhatsappNumber1)
		} else {
			fmt.Printf("\nfailed to sent original message to: %s %s\n", patient.FirstName, patient.LastName)
			err1 := sentMessageTo(patient, patient.WhatsappNumber2)
			if err1 != nil {
				fmt.Printf("\nfailed to sent alternative message to: %s %s\n", patient.FirstName, patient.LastName)
			} else {
				//update message status in db
				fmt.Printf("\nalternative message sent to: %s %s, has this num %s\n", patient.FirstName, patient.LastName, patient.WhatsappNumber2)
			}
		}
	}
}

func sentMessageTo(patient helper.Patient, num string) error {
	message := fmt.Sprintf(
		"السلام عليكم %s %s,\nهذا تذكير بموعد إزالة الدعامة (JJ stent) غداً على الساعة 9 صباحاً.\nالمرجو الحضور إلى العيادة قبل الموعد بساعة واحدة.\nشكراً لتعاونكم.",
		patient.FirstName,
		patient.LastName,
	)

	payload, err := json.Marshal(map[string]string{
		"to":   num,
		"text": message,
	})

	fmt.Println("\npayload: ", string(payload))

	req, err := http.NewRequest("POST", helper.WASENDERA_URL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+helper.WASENDER_TOKEN)
	req.Header.Add("Content-Type", "application/json")

	fmt.Println("\nrequest: ", *req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("\nresponse: ", *resp)
		return errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

package calendar

import (
	"2_18/models"
	"time"
)

type Calendar struct {
	events map[int]map[string][]models.Event
}

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[int]map[string][]models.Event),
	}
}

// CRUD
func (c *Calendar) CreateEvent(userId int, date string, event models.Event) {
	if c.events[userId] == nil {
		c.events[userId] = make(map[string][]models.Event)
	}
	c.events[userId][date] = append(c.events[userId][date], event)
}

func (c *Calendar) GetEventsForDay(userId int, date string) []models.Event {
	if c.events[userId] == nil {
		return []models.Event{}
	}
	return c.events[userId][date]
}

func (c *Calendar) GetEventsForWeek(userId int, date string) ([]models.Event, error) {
	startDay, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	weekday := startDay.Weekday() //вернет текущий день недели
	daysToMonday := int(weekday) - 1
	if daysToMonday < 0 {
		daysToMonday = 6
	}

	monday := startDay.AddDate(0, 0, -daysToMonday)
	weekEvents := []models.Event{}

	for i := 0; i < 7; i++ {
		currentDate := monday.AddDate(0, 0, i)
		dateStr := currentDate.Format("2006-01-02")

		if userEvents, exists := c.events[userId]; exists {
			if dayEvents, exists := userEvents[dateStr]; exists {
				weekEvents = append(weekEvents, dayEvents...)
			}
		}
	}
	return weekEvents, nil
}

func (c *Calendar) GetEventsForMonth(userId int, date string) ([]models.Event, error) {
	startDay, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	firstOfMonth := time.Date(startDay.Year(), startDay.Month(), 1, 0, 0, 0, 0, time.UTC) //первый день нашего месяца
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)                                     //первый день следующего месяца
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)                                     //последний день нашего месяца
	daysInMonth := lastOfMonth.Day()

	monthEvents := []models.Event{}
	for i := 0; i < daysInMonth; i++ {
		currentDay := firstOfMonth.AddDate(0, 0, i)
		dateStr := currentDay.Format("2006-01-02")

		if userEvents, exists := c.events[userId]; exists {
			if dayEvents, exists := userEvents[dateStr]; exists {
				monthEvents = append(monthEvents, dayEvents...)
			}
		}
	}
	return monthEvents, nil
}

func (c *Calendar) DeleteEvent(userId int, date string, description string) bool {
	if c.events[userId] == nil || c.events[userId][date] == nil {
		return false
	}
	dayEvents := c.events[userId][date]
	foundIndex := -1
	for i := range dayEvents {
		if dayEvents[i].Description == description {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return false
	}

	c.events[userId][date] = append(dayEvents[:foundIndex], dayEvents[foundIndex+1:]...)

	if len(c.events[userId][date]) == 0 {
		delete(c.events[userId], date)
	}

	if len(c.events[userId]) == 0 {
		delete(c.events, userId)
	}

	return true
}

func (c *Calendar) UpdateEvent(userId int, oldDate, oldDescription, oldPriority, newDate, newDescription, newPriority string) bool {
	if c.events[userId] == nil || c.events[userId][oldDate] == nil {
		return false
	}

	dayEvents := c.events[userId][oldDate]
	foundIndex := -1
	for i := range dayEvents {
		if dayEvents[i].Description == oldDescription && dayEvents[i].Priority == oldPriority {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return false
	}
	if newDate != "" {
		if newDate != oldDate {
			eventToMove := dayEvents[foundIndex]
			c.events[userId][oldDate] = append(dayEvents[:foundIndex], dayEvents[foundIndex+1:]...)
			eventToMove.Date = newDate

			if c.events[userId][newDate] == nil {
				c.events[userId][newDate] = []models.Event{}
			}
			c.events[userId][newDate] = append(c.events[userId][newDate], eventToMove)
		}
		// По факту тут oldDate == newDate можно ничего не делать так как дата останется, но для понимания оставлю
		//else {
		//	dayEvents[foundIndex].Date = newDate
		//}
	}
	if newDescription != "" {
		dayEvents[foundIndex].Description = newDescription
	}
	if newPriority != "" {
		dayEvents[foundIndex].Priority = newPriority
	}

	return true
}

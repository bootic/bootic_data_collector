package redis_stats

import (
  "datagram.io/data"
  "github.com/vmihailenco/redis"
  "github.com/bitly/go-simplejson"
  "time"
  "fmt"
  "log"
  "strconv"
)

type Tracker struct {
  // Success chan string
  Conn *redis.Client
  Notifier data.EventsChannel
}

func (self *Tracker) Track(key, evtType string) {
  go func(key, evtType string) {
    now := time.Now()

    yearAsString  := strconv.Itoa(now.Year())
    monthAsString := strconv.Itoa(int(now.Month()))
    dayAsString   := strconv.Itoa(now.Day())
    hourAsString  := strconv.Itoa(now.Hour())

    // increment current month in year
    yearKey       := fmt.Sprintf("track:%s:%s:%s", key, evtType, yearAsString)
    self.Conn.HIncrBy(yearKey, monthAsString, 1)

    // increment current day in month
    monthKey      := fmt.Sprintf("track:%s:%s:%s:%s", key, evtType, yearAsString, monthAsString)
    self.Conn.HIncrBy(monthKey, dayAsString, 1)

    // increment current hour in day
    dayKey      := fmt.Sprintf("track:%s:%s:%s:%s:%s", key, evtType, yearAsString, monthAsString, dayAsString)
    self.Conn.HIncrBy(dayKey, hourAsString, 1)

    // Expire day entry after a month
    self.Conn.Expire(dayKey, 2592000)

    // self.Success <- fmt.Sprintf("Done! %s:%s", key, evtType)
  }(key, evtType)
}

func (self *Tracker) Listen() {
  for {
    event := <- self.Notifier
    evtType, _     := event.Get("type").String()
    evtAccount, _  := event.Get("data").Get("account").String()
    self.Track(evtAccount, evtType)
    log.Println("Tracker", evtAccount, " got event", evtType)
  }
}

func (self *Tracker) StoreEvent(event *simplejson.Json) (err error) {
  eventType, err      := event.Get("type").String()
  if err != nil { return }
  eventAccount, err   := event.Get("data").Get("account").String()
  if err != nil { return }
  // Track account
  self.Track(eventType, eventAccount)
  // Track all
  self.Track(eventType, "all")
  return
}

func NewTracker(redisAddress string) (tracker *Tracker, err error) {
  password := "" // no password set
  conn := redis.NewTCPClient(redisAddress, password, -1)
  
  tracker = &Tracker{
    Conn: conn,
    Notifier: make(data.EventsChannel, 1),
  }
  
  go tracker.Listen()
  
  return
}
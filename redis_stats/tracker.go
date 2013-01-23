package redis_stats

import (
  "datagram.io/data"
  "github.com/vmihailenco/redis"
  "time"
  "fmt"
  "strconv"
)

type Tracker struct {
  Conn *redis.Client
  Notifier data.EventsChannel
  Funnels data.EventsChannel
}

func (self *Tracker) TrackTime(accountStr, evtType string) {
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

  }(accountStr, evtType)
}

func (self *Tracker) TrackFunnel(accountStr, evtType, statusStr string) {
  go func(key, evtType, statusStr string) {
    now := time.Now()

    yearAsString  := strconv.Itoa(now.Year())
    monthAsString := strconv.Itoa(int(now.Month()))

    // increment current month in year
    yearKey       := fmt.Sprintf("funnels:%s:%s:%s", key, evtType, yearAsString)
    self.Conn.HIncrBy(yearKey, statusStr, 1)

    // increment current day in month
    monthKey      := fmt.Sprintf("funnels:%s:%s:%s:%s", key, evtType, yearAsString, monthAsString)
    self.Conn.HIncrBy(monthKey, statusStr, 1)

  }(accountStr, evtType, statusStr)
}

func (self *Tracker) listenForPageviews() {
  for {
    event := <- self.Notifier
    evtType, _     := event.Get("type").String()
    evtAccount, _  := event.Get("data").Get("account").String()
    self.TrackTime(evtAccount, evtType)
    self.TrackTime("all", evtType)
  }
}

func (self *Tracker) listenForFunnels() {
  for {
    event := <- self.Funnels
    evtType, _     := event.Get("type").String()
    evtAccount, _  := event.Get("data").Get("account").String()
    evtStatus, _  := event.Get("data").Get("status").String()
    self.TrackFunnel(evtAccount, evtType, evtStatus)
    self.TrackFunnel("all", evtType, evtStatus)
  }
}

func NewTracker(redisAddress string) (tracker *Tracker, err error) {
  password := "" // no password set
  conn := redis.NewTCPClient(redisAddress, password, -1)
  
  tracker = &Tracker{
    Conn: conn,
    Notifier: make(data.EventsChannel, 1),
    Funnels: make(data.EventsChannel, 1),
  }
  
  go tracker.listenForPageviews()
  go tracker.listenForFunnels()
  
  return
}
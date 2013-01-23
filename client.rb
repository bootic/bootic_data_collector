require 'socket'
require 'json'
s = UDPSocket.new

# HOST = '192.168.245'
HOST = 'localhost'

N = 70
ACCOUNTS = ['romano', 'japijane', 'depto51']
EVENTS = ['pageview', 'request', 'apievent', 'order']
STATUSES = ['draft', 'draft', 'draft', 'checkout', 'closed']

# app, account/shop, user/token
1.upto N do |i|
  msg = JSON.dump(
    type: EVENTS[rand(EVENTS.size)],
    time: Time.now.to_s,
    data: {
      app: "app-#{i}",
      account: ACCOUNTS[rand(ACCOUNTS.size)],
      status: STATUSES[rand(STATUSES.size)],
      user: "user-#{i}",
      resource_id: i,
      browser: {
        os: 'Mac',
        name: 'Chrome',
        version: 23
      }
    }
  )
  s.send(msg, 0, HOST, 5555)  
  p [:send, msg]
  sleep 0.1
end

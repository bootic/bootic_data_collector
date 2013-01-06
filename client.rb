require 'socket'
require 'json'
s = UDPSocket.new

# HOST = '192.168.245'
HOST = 'localhost'

N = 70
# app, account/shop, user/token
1.upto N do |i|
  msg = JSON.dump(
    type: "event-##{i}",
    time: Time.now.to_s,
    data: {
      app: "app-#{i}",
      account: "account-#{i}",
      status: "status-#{i}",
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

require 'socket'
require 'json'
s = UDPSocket.new

# HOST = '192.168.245'
HOST = 'localhost'

1.upto 10 do |i|
  msg = JSON.dump(
    desc: "hello ##{i}",
    tags: ['tag1', 'tag2', 'tag3', "tag-x-#{i}"]
  )
  s.send(msg, 0, HOST, 5555)  
  p [:send, msg]
  sleep 0.5
end

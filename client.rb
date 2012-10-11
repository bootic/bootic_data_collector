require 'socket'
s = UDPSocket.new

# HOST = '192.168.245'
HOST = '192.168.100.245'

1.upto 1000 do |i|
  msg = "hello ##{i}"
  s.send(msg, 0, HOST, 5555)  
  p [:send, msg]
  sleep 0.5
end

# converts .env file to docker run args (ie: `-e X=Y`)
# If you pass in the --type flag, you can change it to output Powershell (ps) or bash

require 'optparse'

body = File.read(ARGV[0])

options = {}
OptionParser.new do |opt|
  opt.on('--type TYPE') { |o| options[:type] = o }
end.parse!

output = " "
body.each_line do |line|
#   print "#{line}"
  line = line.strip
  next if line == ""
  next if line[0] == '#'
  case options[:type]
  when 'ps'
    split = line.split('=', 2)
    output += " $env:#{split[0]} = \"#{split[1]}\";"
  else
    output += " -e \"#{line}\""
  end
end

puts output

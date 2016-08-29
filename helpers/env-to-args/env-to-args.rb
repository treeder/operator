
# converts .env file to docker run args (ie: `-e X=Y`)
body = File.read(ARGV[0])

output = " "
body.each_line do |line|
#   print "#{line}"
  line = line.strip
  next if line == ""
  output += " -e \"#{line}\""
end

puts output

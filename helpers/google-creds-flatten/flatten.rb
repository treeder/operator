require 'base64'

# converts .env file to docker run args (ie: `-e X=Y`)
body = File.read(ARGV[0])

output = ""
body.each_line do |line|
#   print "#{line}"
  line = line.strip
  next if line == ""
  # output += line.gsub(/\s+/, "")
  output += line.gsub("\n", "")
end
output = Base64.strict_encode64(output)
puts output

# converts an RSA key to a single line with \n in between, for use in .env files

body = File.read(ARGV[0])

output = ""
body.each_line do |line|
  # print "#{line}"
  line = line.strip
  next if line == ""
  output += "#{line}\\n"
end

puts output

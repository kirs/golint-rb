require "golint/version"
require "tempfile"
require "base64"

module Golint
  class Match < Struct.new(:line, :comment)
  end

  class << self
    REGEXP = /([\d]*):([\d]*):([\s\S]*)$/

    def lint(content)
      code = Base64.encode64(content).strip[0..16]
      file = Tempfile.new(code)
      file.write(content)
      file.close

      diff = `golint #{file.path}`

      pattern = "#{file.path}:"

      matches = []
      diff.each_line do |line|
        line = line.sub(pattern, "").sub("\n", "")
        res = line.match(REGEXP)
        matches << Match.new(res[1], res[3].strip)
      end

      matches
    end
  end
end

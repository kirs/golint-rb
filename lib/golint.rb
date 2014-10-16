require "golint/version"
require "tempfile"
require "base64"
require 'open3'

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

      stdin, stdout, stderr, wait_thr = Open3.popen3('golint', file.path)

      out = stdout.read
      err = stderr.read

      if err.size > 0
        parse_matches(file.path, err)
      elsif out.size > 0
        parse_matches(file.path, out)
      else
        []
      end
    end

    def parse_matches(path, body)
      pattern = "#{path}:"
      matches = []
      body.each_line do |line|
        line = line.sub(pattern, "").sub("\n", "")

        res = line.match(REGEXP)
        matches << Match.new(res[1].to_i, res[3].strip)
      end

      matches
    end
  end
end

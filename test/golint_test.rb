require 'test_helper'

class TestGolint < MiniTest::Unit::TestCase
  def test_shit
    file = File.open('test/fixtures/server.go')
    puts Golint.lint(file.read)
  end
end

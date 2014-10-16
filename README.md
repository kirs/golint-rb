# Golint-rb

[![Build Status](https://travis-ci.org/kirs/golint-rb.svg)](https://travis-ci.org/kirs/golint-rb)

Lints your Go code with golint utility.

## Installation

First of all, get `golint` package:

```
go get github.com/golang/lint/golint
```

Add this line to your application's Gemfile:

    gem 'golint'

And then execute:

    $ bundle

Or install it yourself as:

    $ gem install golint

## Usage

```ruby
matches = Golint.lint("some of your go code")
matches.each do |m|
  puts m.line
  puts m.comment
end
```

``

## Contributing

1. Fork it ( https://github.com/[my-github-username]/gofmt-rb/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

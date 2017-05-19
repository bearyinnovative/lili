#!/usr/bin/env ruby

require 'json'

output = `./_internal/bearychat-zhihu-v3.sh`
data = JSON.parse(output)

data.each do |d|
  title = d["title"].gsub(/<em>/, '').gsub(/<\/em>/, '')

  if d['question'] != nil and d['question'] != nil and d["author"] != nil
    url = "https://www.zhihu.com/question/#{d['question']}/answer/#{d['answer']}"
    puts "* Q: #{title}\n  * A: #{d["author"]}: #{url}"
  else
    puts "* #{title}"
  end
end
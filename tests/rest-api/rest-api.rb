#!/usr/bin/env ruby
#
# Copyright 2014 Rohith All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
require 'optparse'
require 'sinatra'
require 'yaml'
require 'rack'

module RestAPI
  class Application < Sinatra::Base
    enable :logging, :static, :raise_errors
    set :filename, "methods.yml"

    unless File.exists? settings.filename
      raise ArgumentError, "you need to specify create #{settings.filename} rest methods"
    end

    YAML.load(File.read(settings.filename)).each do |rest|
      method_type = rest['method'].downcase
      next unless method_type =~ /^(post|get|delete|put)$/
      send(method_type,rest['uri']) do
        write rest['content']
      end
    end

    def write message
      content_type :json
      message
    end
  end
end

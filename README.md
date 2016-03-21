LruCache [![Build Status](https://travis-ci.org/die-net/lrucache.svg?branch=master)](https://travis-ci.org/die-net/lrucache)
========

LruCache is a thread-safe, in-memory [httpcache.Cache](https://github.com/gregjones/httpcache) implementation that evicts the least recently used entries when a byte size limit would be exceeded.

See the [godoc documentation](https://godoc.org/github.com/die-net/lrucache) for more detail.

Also included are a test suite with close to 100% test coverage and a parallel benchmark suite that can be used to verify thread-safely. Benchmarks show individual Set, Get, and Delete operations take under 400ns to complete.

License
-------

Copyright 2016 Aaron Hopkins and contributors

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at: http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

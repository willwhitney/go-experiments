## crawler

`crawler.go` is a small chunk of Go that randomly crawls through Wikipedia with 200 parallel goroutines, printing the titles of the pages it reaches to stdout.

## server

`server.go` is a quick performance test of Go that factors a thousand numbers in parallel using the least efficient algorithm.
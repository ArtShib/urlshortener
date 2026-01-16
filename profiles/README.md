File: audit.test
Build ID: 318d062002c3fd67c9b005e9941177d01fa7dadf
Type: alloc_space
Time: 2025-12-29 19:55:13 MSK
Showing nodes accounting for -328.33MB, 90.63% of 362.27MB total
Dropped 32 nodes (cum <= 1.81MB)
flat  flat%   sum%        cum   cum%
-203.02MB 56.04% 56.04%  -349.53MB 96.48%  context.WithDeadlineCause
-146.52MB 40.44% 96.48%  -146.52MB 40.44%  time.newTimer
21.70MB  5.99% 90.49%    21.70MB  5.99%  github.com/ArtShib/urlshortener/internal/workerpool/audit.New (inline)
-0.50MB  0.14% 90.63%  -350.03MB 96.62%  github.com/ArtShib/urlshortener/internal/workerpool/audit.(*WorkerPoolEvent).worker
0     0% 90.63%  -349.53MB 96.48%  context.WithDeadline (inline)
0     0% 90.63%  -349.53MB 96.48%  context.WithTimeout
0     0% 90.63%  -349.53MB 96.48%  github.com/ArtShib/urlshortener/internal/workerpool/audit.(*WorkerPoolEvent).processEvent
0     0% 90.63%    21.70MB  5.99%  github.com/ArtShib/urlshortener/internal/workerpool/audit.BenchmarkWorkerPool_Base
0     0% 90.63%    21.70MB  5.99%  testing.(*B).launch
0     0% 90.63%    21.69MB  5.99%  testing.(*B).runN
0     0% 90.63%  -146.52MB 40.44%  time.AfterFunc
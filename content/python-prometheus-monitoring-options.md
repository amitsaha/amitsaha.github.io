Title:Options for Monitoring Python applications with Prometheus
Date: 2018-01-23 13:00
Category: Python
Status: Draft


In an earlier article, [Monitoring Your Synchronous Python Web Applications Using Prometheus](https://blog.codeship.com/monitoring-your-synchronous-python-web-applications-using-prometheus/), I discussed a limitation of using the Python client for prometheus. 

##  Limitation of native prometheus exporting

[prometheus](https://prometheus.io) was built with single process multi-threaded applications in mind.
I use the term multi-threaded here to also include coroutine based concurrent applications such as
those written in `golang` or using Python's asynchronous primitives 
(Example: [Monitoring Your Asynchronous Python Web Applications Using Prometheus](https://blog.codeship.com/monitoring-your-asynchronous-python-web-applications-using-prometheus/)). 

The problem arises when you have multi-process applications using the model of `worker processes such as:

- WSGI web applications deployed  via `uwsgi`or `gunicorn`
- Thrift RPC servers

When you have these multi-process applications (as a single application instance), 
you get into a situation where any of the multiple workers
can respond to prometheus's scraping request and hence you can have one scrape response having a value
of a `counter` as `200` and the immediate next scrape having a counter value of `100`. To prometheus, they
seem like actual values since they are the same metric. The same inconsistent behaviour can happen with a
`gauge` or a `histogram`.

To work around this problem, you can add a unique `worker_id` as a label such that each metric as scraped
by prometheus is unique for one application instance (by virtue of having different value for the 
label, `worker_id`). However, even this introduces the same inconsistent behavior as above at least
for the short term only for `counter` metrics.  To explain what I mean - let's say we have a counter metric 
`request_count`. This is the number of requests served by the web application. Now consider, two time windows: 
`t1-t2` and `t2-t3`.  Ideally, if we have N workers, requests to web application will
be served such that requests over a time window are equally distributed
among the N workers. However, there's nothing preventing a
single worker serving all 100 requests in the time window `t1-t2`
and a different worker serving all 50 requests in the time window
`t2-t3`. If you now do `sum(request_count) by (instance)` (for example)
you will see the value of the counter as decreasing. This isn't a problem
with `gauges` or `histogram` metric types.

In addition, this leads to a proliferation of metrics: 
for a single metric we now have `# of workers x metric` number of 
metrics per application instance. 

## Multi-process mode

To solve the above issues, the prometheus [Python Client](https://github.com/prometheus/client_python)
has a multi-processing mode which essentially creates a shared prometheus registry and shares
it among all the processes. 


The demos can be found in the [python-prometheus-demo](https://github.com/amitsaha/python-prometheus-demo) repository.

## StatsD Exporter

## Conclusion




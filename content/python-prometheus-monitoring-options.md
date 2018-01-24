Title: Your options for monitoring multi-process Python applications with Prometheus
Date: 2018-01-23 13:00
Category: Python
Status: Draft


In an earlier article, [Monitoring Your Synchronous Python Web Applications Using Prometheus](https://blog.codeship.com/monitoring-your-synchronous-python-web-applications-using-prometheus/), I discussed a limitation of using the Python client for prometheus. 

##  Limitation of native prometheus exporting

[prometheus](https://prometheus.io) was built with single process multi-threaded applications in mind.
I use the term multi-threaded here to also include coroutine based concurrent applications such as
those written in `golang` or using Python's asynchronous primitives 
(Example: [Monitoring Your Asynchronous Python Web Applications Using Prometheus](https://blog.codeship.com/monitoring-your-asynchronous-python-web-applications-using-prometheus/)). 

The problem arises when you have multi-process applications using the model of independent
processes such as WSGI web applications deployed  via `uwsgi`or `gunicorn`.

When you have these multi-process applications (as a single application instance), 
you get into a situation where any of the multiple workers
can respond to prometheus's scraping request and hence you can have one scrape response having a value
of a `counter` as `200` and the immediate next scrape having a counter value of `100`. To prometheus, they
seem like actual values since they are the same metric. The same inconsistent behaviour can happen with a
`gauge` or a `histogram`. You really want each scraping response to return the `overall` values for a
metric rather than what each invidividual worker thinks (a.k.a [WYSIATI](https://jeffreysaltzman.wordpress.com/2013/04/08/wysiati/))
the value to be.

## Solution #1 - Add a unique label to each metric

To work around this problem, you can add a unique `worker_id` as a label such that each metric as scraped
by prometheus is unique for one application instance (by virtue of having different value for the 
label, `worker_id`). Then we can perform aggregation such as:

```
sum by (instance, http_status) (sum without (worker_id) (rate(request_count[5m])))
```

This will perform the aggregation across all `worker_id` metrics (basically ignoring it)
and then we can group by the `instance` and any other label associated to the metric
of interest.

One point worth noting here is that this leads to a proliferation of metrics: 
for a single metric we now have `# of workers x metric` number of 
metrics per application instance. 

A demo of this approach can be found [here](https://github.com/amitsaha/python-prometheus-demo/tree/master/flask_app_prometheus_worker_id).

## Solution #2: Multi-process mode

The prometheus [Python Client](https://github.com/prometheus/client_python)
has a multi-processing mode which essentially creates a shared prometheus registry and shares
it among all the processes and hence the [aggregation](https://github.com/prometheus/client_python/blob/master/prometheus_client/multiprocess.py
) happens at the application level. When,
prometheus scrapes the application instance, no matter which worker responds to the scraping
request, the metrics reported back describes the application's behaviour, rather than 
the worker responding.

A demo of this approach can be found [here](https://github.com/amitsaha/python-prometheus-demo/tree/master/flask_app_prometheus_multiprocessing).

## Solution #3: The Django way

## Solution #4: StatsD exporter

I discussed this solution in [Monitoring Your Synchronous Python Web Applications Using Prometheus](https://blog.codeship.com/monitoring-your-synchronous-python-web-applications-using-prometheus/). Essentially, instead of exporting native prometheus metrics from your
application and prometheus scraping our application, we push our metrics to a locally running [statsd exporter](https://github.com/prometheus/statsd_exporter) instance. Then, we setup prometheus to scrape the statsd exporter instance.


## Exporting metrics for non-HTTP applications



## Resources

- [Common query patterns in PromQL](https://www.robustperception.io/common-query-patterns-in-promql/)




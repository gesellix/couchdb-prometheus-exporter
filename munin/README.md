# Munin Plugin for CouchDB #

[munin-plugin-couchdb][0] is the [Munin][1] plugin that allows to monitor
[Apache CouchDB][2] instance.


## Install and Setup ##

First of all, ensure that your system has installed [Perl 5.12+][3] and
two additional libraries: [LWP::UserAgent][4] and [JSON][5]. Sure, you would
also need to have [Munin][1] installed.

The plugin installation is quite trivial operation

1. `git clone https://github.com/gws/munin-plugin-couchdb`
2. `cp munin-plugin/couchdb_ /etc/munin/plugins`
3. Read the docs via `perldoc couchdb_`
4. Create `/etc/munin/plugin-conf.d/couchdb` and setup proper configuration
5. Check that it works: `su munin -c '/usr/sbin/munin-run couchdb_'`

For additional information about plugins installation consult with
[Munin docs][6].


### Setting Auth Credentials ###

The `munin-plugin-couchdb` is able to gather statistics not only
from [/_stats][13] resource (which doesn't requires any authentication by
default unless [require_valid_user][14] is set on) but also from other
resources like [/_active_tasks][15] which requires to provide CouchDB server
administrator's credentials.

Leaving such credentials in plain text within config file is **dangerous**, so
make sure that plugin's configuration file is readable only for trusted users.

If you're going to monitor remote server (not on localhost) make sure that
you're using secure connection with it (HTTPS or SSH-tunnel) to not transfer
credentials in plain text over the network.


## Monitoring ##


### HTTPD Metrics ###

These metrics are related to [Mochiweb](https://github.com/mochi/mochiweb) -
the CouchDB's HTTP server which runs the API and communicates with the world.

This plugin also gathers all metrics from CouchDB via HTTP API, so it causes
so overhead: one request to fetch `max_dbs_open` from [/_config][16] resource,
one request to fetch all stats from [/_stats][13], 3 more requests (by default)
for `couchdb_request_times` graph per each sample and optional request to
[/_active_tasks][15] if allowed plus one per each monitored database. In total
at least 6 requests per stats update.


#### Request Methods ####

The `couchdb_httpd_request_methods` graph provides information about all HTTP
requests in context of used method. It counts the next methods:

- `HEAD`
- `GET`
- `POST`
- `PUT`
- `DELETE`
- `COPY`

![Request methods](http://gws.github.io/munin-plugin-couchdb/images/request-methods.png)


#### Request Times ####

The `couchdb_request_times` graph shows stddev/mean of HTTP request time within
each [sampling range][17].

In CouchDB configuration information isn't available, the default samples
(`[60, 300, 900]`) will be used. Note, that in this case for each sample that
doesn't match value defined in [stats/samples][17] option this graph will
print zeros.

![Request times](http://gws.github.io/munin-plugin-couchdb/images/request-times.png)


#### Requests by Type ####

The `couchdb_httpd_requests` graph shows rate of HTTP requests in context of
their type:

- `HTTP requests`: overall amount of HTTP requests
- `bulk requests`: how often were used [bulk updates][8]
- `view reads`: amount of requests to the [view indexes][9]
- `temporary view reads`: amount of requests to the [temporary view indexes][10]

![Requests by type](http://gws.github.io/munin-plugin-couchdb/images/http-requests.png)


#### Continuous Changes Feeds Listeners ####

While `clients_requesting_changes` metric is in the same group as
`bulk_requests`, `temporary_view_reads` and others,
the `couchdb_clients_requesting_changes` graph shows not requests rate, but
the current amount of active clients to continuous changes feeds.

This graph also helps to roughly estimate amount of continuous replications
that are running against monitored instance.

![Continuous changes feeds listeners](http://gws.github.io/munin-plugin-couchdb/images/continuous-changes-feeds-listeners.png)


#### Response Status Codes ####

The `couchdb_httpd_status_codes` graph provides information about HTTP
responses in context of [status code][7].

Keeping eye on amount of HTTP `4xx` and `5xx` responses helps you provide
quality service for you users. Normally, you want to see no `500` errors at all.
Having high amount of `401` errors could say about authentication problems
while `403` tell you that something or someone actively doing things that he's
shouldn't do.

![Response status codes](http://gws.github.io/munin-plugin-couchdb/images/http-status-codes.png)


### Server Metrics ###

These metrics are related to whole server instance.


#### Authentication Cache ####

The `couchdb_auth_cache` graph shows rate of authentication cache hits/misses.

CouchDB keeps some amount of user credentials in memory to speedup
authentication process by elimination of additional database lookups.
This cache size is limited by the configuration option [auth_cache_size][11].
On what this affects? In short, when user login CouchDB first looks for user
credentials what associated with provided login name in auth cache and if they
miss there then it reads credentials from auth database (in other words,
from disk).

The `auth_cache_miss` metric is highly related to `HTTP 401` responses one,
so there are three cases that are worth to be looked for:

- High `cache misses` and high `401` responses: something brute forces your
  server by iterating over large set of user names that doesn't exists for your
  instance

- High `cache misses` and low `401` responses: your `auth_cache size` is
  too small to handle all your active users. try to increase his capacity
  to reduce disk I/O

- Low `cache misses` and high `401` responses: much likely something tries
  to brute force passwords for existed accounts on your server

Note that "high" and "low" in metrics world should be read as "anomaly high"
and "anomaly low".

Ok, but why do we need auth cache hit then? We need it as an ideal value
to compare misses counter with. Just for instance, is 10 cache misses a high
value? What about 100 or 1000? Having cache hits rate at some point helps
to answer on this question.

![Authentication cache ratio](http://gws.github.io/munin-plugin-couchdb/images/auth-cache.png "Whoa! Someone really tries to do bad things there")


#### Databases I/O ####

The `couchdb_database_io` graph shows overall databases read/write rate.

![Databases I/O](http://gws.github.io/munin-plugin-couchdb/images/database-io.png)


#### Open Databases ####

The `couchdb_open_databases` graph shows amount of currently opened databases.

CouchDB only keeps opened databases which are receives some activity: been
requested or running the compaction. The maximum amount of opened
databases in the same moment of time is limited by [max_dbs_open][12]
configuration option.  When CouchDB hits this limit, any request to "closed"
databases will generate the error response: `{error, all_dbs_active}`.

However, once opened database doesn't remains open forever: in case of
inactivity CouchDB eventually closes it providing more space in the room for
others, but sometimes such cleanup may not help. This graph's goal is to help
you setup correct `max_dbs_open` value that'll fit your needs.

*Notice:* If server administrator's credentials provided (need to request
[/_config][16] resource) the `max_dbs_open` configuration value will be used to
set proper `warning` and `critical` levels.

![Open databases](http://gws.github.io/munin-plugin-couchdb/images/open-databases.png)


#### Open Files ####

The `couchdb_open_files` graph shows amount of currently opened file
descriptors.

*Notice:* Handling system `nofile` limit isn't implemented yet and couldn't be
possible for remote instances.

![Open files](http://gws.github.io/munin-plugin-couchdb/images/open-files.png)


#### Active Tasks ####

**Warning:** this graph is *disabled* by default. To enable it you should
set `env.monitor_active_tasks yes` in plugin configuration file and also
provide CouchDB server administrator user. See `Setting Auth Credentials`
section above for recommendations.


The `couchdb_active_tasks` graph shows current processes that runs on CouchDB:

- Active replications, served by this CouchDB instance
- View index builds
- Database and views compactions

This information is very valuable since some of these operations are very IO
heavy (compactions are so). For instance, you're looking on `diskstats_iops`
graph and see high write activity, but for most cases you could say for sure
who generates it. Combining these graphs together for the same period may
give you the answer is this activity is related to CouchDB and how if it is.

![Active tasks](http://gws.github.io/munin-plugin-couchdb/images/active-tasks.png)


#### Users ####

**Warning:** these graphs are *disabled* by default. To enable them you should
set `env.monitor_users yes` in plugin configuration file and also
provide CouchDB server administrator user. See `Setting Auth Credentials`
section above for recommendations.


The `couchdb_users` and `couchdb_admin_users` graphs shows total amount of known
users by CouchDB.

The `couchdb_admin_users` graph is stand alone to easily track amount of server
administrators. In most time their number is stable and any unexpectable changes
may be a sign for worry about server security.

The `couchdb_users` graph shows users from [authentication database][18] and
tracks `registered` and `deleted` amount of them. This helps to estimate size of
your users database growing and decreasing in time.

![CouchDB users](http://gws.github.io/munin-plugin-couchdb/images/users.png)


### Database Metrics ###

[munin-plugin-couchdb][0] also allows to monitor few databases metrics that could
be useful. To enable it you need to set `env.monitor_databases yes` variable
in your plugin's configuration file and explicitly define list of databases
which would be monitored in `env.databases`. For example:

    [couchdb]
    env.uri    http://localhost:5984
    env.username  admin
    env.password  s3cR1t
    env.monitor_databases  yes
    env.databases  mailbox, db/with/slashes, data+ba$ed

Note, that user for provided credential should have read access to the specified
databases to request [database information][19] from them.


#### Documents Count ####

The `couchdb_db_${dbname}_docs` graph shows amount of existed and deleted
documents in specific database.

CouchDB doesn't physically removes documents on `DELETE` leaving tombstone
instead to be able replicate this information to others databases and to prevent
accidental  "resurrection" of such documents during push replication.

However, when amount of deleted documents becomes significantly greater than
existed ones, this may seriously affect on consumed disk space. Such "graveyard
databases" are needed in cleanup from deleted documents (in case when it's ever
possible) and this graph helps to detect them.

![Database documents](http://gws.github.io/munin-plugin-couchdb/images/database-docs.png)


#### Database Fragmentation ####

The `couchdb_db_${dbname}_frag` graph tracks database `disk_size` grow in time
and overhead caused over `data_size`.

Databases are needs to be compacted from time to time to retain used disk space
by old documents revisions, but it's hard to note when compaction is worth to
run especially since it's heavy disk I/O operation: you probably wouldn't
compact 1TiB database just to free 20GiB. This graph helps to find answers on
these two questions: "when?" and "how much?".

![Database disk usage](http://gws.github.io/munin-plugin-couchdb/images/database-frag.png)


## License ##

[Beerware](https://tldrlegal.com/license/beerware-license)


[0]: https://github.com/gws/munin-plugin-couchdb
[1]: http://munin-monitoring.org/
[2]: http://couchdb.apache.org/
[3]: http://www.perl.org/
[4]: http://search.cpan.org/dist/LWP-UserAgent-Determined/
[5]: http://search.cpan.org/dist/JSON/
[6]: https://munin.readthedocs.org/en/latest/plugin/use.html#installing
[7]: http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
[8]: http://docs.couchdb.org/en/latest/api/database/bulk-api.html#post--db-_bulk_docs
[9]: http://docs.couchdb.org/en/latest/api/ddoc/views.html
[10]: http://docs.couchdb.org/en/latest/api/database/temp-views.html#post--db-_temp_view
[11]: http://docs.couchdb.org/en/latest/config/auth.html#couch_httpd_auth/auth_cache_size
[12]: http://docs.couchdb.org/en/latest/config/couchdb.html#couchdb/max_dbs_open
[13]: http://docs.couchdb.org/en/latest/api/server/common.html#stats
[14]: http://docs.couchdb.org/en/latest/config/auth.html#couch_httpd_auth/require_valid_user
[15]: http://docs.couchdb.org/en/latest/api/server/common.html#active-tasks
[16]: http://docs.couchdb.org/en/latest/api/server/configuration.html#get--_config
[17]: http://docs.couchdb.org/en/latest/config/misc.html#stats/samples
[18]: http://docs.couchdb.org/en/latest/intro/security.html#authentication-database
[19]: http://docs.couchdb.org/en/latest/api/database/common.html#get--db

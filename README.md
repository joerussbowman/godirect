godirect
========

godirect is a tool for managing http redirects for
multiple hosts.

Note: godirect is still a work in progress. The following
TODO items need to be completed.

TODO:
 - godirect config management tool
 - godirect engine config management listener and handlers
 - finalize config structure for godirect engine
 - provide better support for request types other than GET.
   ie: a POST, regardless of path, should always proxy
 - benchmarks to validate performance is acceptable

This tool is to solve the case where organizations manage
http redirects via virtualhosts in tools like apache. Apache
can handle this using RewriteMap and dbm hash files, but
the business process for getting these in place still requires
some work. Adding a new hostname can still require a restart
of the apache processes.

godirect is a two part tool to make this easier. The first
tool is an engine which can be dynamically configured to
add new host names and redirects for them. It works simply by
being placed in front of you normal web server. When a request
is made it matches the host name to a configured host and compares
the path of the request to a list of redirects. If the path is
in the redirect list it will redirect the browser. Otherwise, godirect
will proxy the request to the backend web server.

The second tool (still not built) is a web interface used to manage
godirect engine servers. The godirect engine supports getting updated
configuration via JSON over an internal http listener. The godirect
config manager allows you to view and modify the configuration and
push changes to multiple godirect engines. This removes the need
for operation teams to be involved in managing redirect changes.


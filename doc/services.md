Services are the basic building blocks of Sigil. They provide useful functionality that is either
core (HTTP server, Websocket server) or modular/interchangable (storage backends, e.g. MySQL, image
manipulation, template compiler etc.), and can bind themselves to present select functionality to
the Javascript environment.

Requirements:
  * Methods for Setup, Start, and Stop operations. Services should be able to cleanly cease
    operations, with a possibility of no-downtime restarts.
# forum
a bare bones forum

This "forum" was hacked out as quickly as possible, with an aim towards using as
few "external" dependencies as possible. Normally I would use things like sqlx,
gorilla mux (and other gorilla bits), etc. The only item included from
"external" sources is the sqlite database connector.

Also, since it's done as quickly as possible, there are no tests, to speak of.
There is one small set of tests for the user persistence, but I did that only to
confirm that I basically had it working before doing the rest of the model.

There is no authentication, this just accepts a new user name of whoever wanders
by.
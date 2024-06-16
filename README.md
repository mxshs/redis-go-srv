# TODOS

* Come up with a better structure. Perhaps passing connections around is a bad idea.
* Fix the parser. Right now, since I wrote a recursive descent one, I have to add dubious checks and offset shifts to account for dumb decisions.
* Rework replication. Creating a separate goroutine for slave <-> master communication before running the actual handler looks ugly.

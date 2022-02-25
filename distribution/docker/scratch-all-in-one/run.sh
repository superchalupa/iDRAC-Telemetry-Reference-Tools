#!/bin/sh
PATH=$PATH:/app/.

echo STARTING DAEMONS...

dbdiscauth &
configui &
redfishread &
dbdiscauth


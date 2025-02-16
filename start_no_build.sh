#!/bin/bash

# Define the path to the PID file
PIDFILE="/home/markort/var/run/bekindfrontend.pid"

# Check if the PID file exists
if [ -f $PIDFILE ]; then
  # Read the PID from the file
  PID=$(cat $PIDFILE)
  
  # Check if the process is running
  if ps -p $PID > /dev/null; then
    echo "BeKindFrontend is already running with PID $PID"
    exit 1
  else
    echo "PID file found but process not running. Removing stale PID file."
    rm $PIDFILE
  fi
fi

# Run the application using nohup and store the PID
nohup go run ./cmd/. > /dev/null 2>&1 &
echo $! > $PIDFILE

echo "BeKindFrontend started with PID $(cat $PIDFILE)"
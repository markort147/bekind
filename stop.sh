#!/bin/bash

# Define the path to the PID file
PIDFILE="/home/markort/var/run/bekindfrontend.pid"

# Check if the PID file exists
if [ -f $PIDFILE ]; then
  # Read the PID from the file
  PID=$(cat $PIDFILE)
  
  # Find child PIDs of the parent process
  CHILD_PIDS=$(pgrep -P $PID)
  
  # Kill the child processes
  if [ -n "$CHILD_PIDS" ]; then
    kill $CHILD_PIDS
  fi

  # Kill the process
  kill $PID
  
  # Remove the PID file
  rm $PIDFILE
  
  echo "BeKindFrontend stopped"
else
  echo "PID file not found. Is BeKindFrontend running?"
fi
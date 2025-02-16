#!/bin/bash

# Check if the version argument is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1
APP_PATH="./app_v$VERSION"

# Check if the application file exists
if [ ! -f $APP_PATH ]; then
  echo "Application file $APP_PATH does not exist"
  exit 1
fi

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
nohup $APP_PATH > /dev/null 2>&1 &
echo $! > $PIDFILE

echo "BeKindFrontend started with PID $(cat $PIDFILE)"
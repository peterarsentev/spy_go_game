#!/bin/bash

APP_NAME="spy.1.0.1"
LOG_FILE="spy.log"
PID_FILE="spy.pid"

build() {
  echo "Building app for Linux..."
  export GOOS=linux
  export GOARCH=amd64
  go build -o $APP_NAME ./cmd/main.go
}

start() {
  if [ -f $PID_FILE ] && kill -0 $(cat $PID_FILE) 2>/dev/null; then
    echo "$APP_NAME is already running with PID $(cat $PID_FILE)"
    exit 1
  fi
  echo "Starting $APP_NAME..."
  nohup ./$APP_NAME > $LOG_FILE 2>&1 &
  echo $! > $PID_FILE
  echo "$APP_NAME started with PID $(cat $PID_FILE)"
}

stop() {
  if [ ! -f $PID_FILE ]; then
    echo "No PID file found. Is $APP_NAME running?"
    exit 1
  fi
  PID=$(cat $PID_FILE)
  echo "Stopping $APP_NAME with PID $PID..."
  kill $PID
  rm -f $PID_FILE
  echo "Stopped."
}

status() {
  if [ -f $PID_FILE ] && kill -0 $(cat $PID_FILE) 2>/dev/null; then
    echo "$APP_NAME is running with PID $(cat $PID_FILE)"
  else
    echo "$APP_NAME is not running"
  fi
}

case "$1" in
  build)
    build
    ;;
  start)
    start
    ;;
  stop)
    stop
    ;;
  status)
    status
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {build|start|stop|status|restart}"
    exit 1
    ;;
esac

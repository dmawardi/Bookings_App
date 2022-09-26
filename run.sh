#!/bin/bash

# go build cmd/web AND output to bookings
go build -o bookings cmd/web && ./bookings
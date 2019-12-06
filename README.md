# golang-mqtt-stats

# MQTT Topic Statistics

## Overview

This is a simple MQTT subscriber client in Go to display topic performance statistics.
It is the command-line interface equivalent of https://github.com/gambitcomminc/mqtt-stats


## Installation / Requirements

* Follow instructions for Eclipse Paho Go library https://github.com/eclipse/paho.mqtt.golang

* go build golang-mqtt-stats.go

## Usage

This shows with 10 messages / second from the MIMIC MQTT Simulator Bosch simulation

     % ./golang-mqtt-stats --host 192.9.192.247 --topic '#'
     Subscribing to topic # on 192.9.192.247:1883
     elapsed 5.540: msgs/s 9.4, bytes/s 3045.6, bytes/msg 324.5
     elapsed 6.003: msgs/s 10.0, bytes/s 3305.0, bytes/msg 330.7
     elapsed 5.005: msgs/s 10.0, bytes/s 3302.5, bytes/msg 330.6
     elapsed 5.003: msgs/s 10.0, bytes/s 3304.4, bytes/msg 330.6
     elapsed 5.005: msgs/s 10.0, bytes/s 3304.3, bytes/msg 330.8


#!/bin/sh

curl -X DELETE http://localhost:8083/connectors/inventory-connector

sleep 2

curl -i -X POST -H "Accept:application/json" -H  "Content-Type:application/json" http://127.0.0.1:8083/connectors/ -d @register-mongo.json

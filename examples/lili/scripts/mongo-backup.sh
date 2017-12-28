#!/bin/bash

mongodump -v --db lili --host $MONGO_SERVER --out=/tmp/mongodump/
mongodump -v --db house --host $MONGO_SERVER --out=/tmp/mongodump/
#!/bin/bash

mongodump -v --db lili --host $MONGO_SERVER --out=/tmp/mongodump/lili
mongodump -v --db house --host $MONGO_SERVER --out=/tmp/mongodump/house
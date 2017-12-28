#!/bin/bash

cd /tmp/mongodump

dirs=(*)

for dir in "${dirs[@]}"
do
	echo "mongorestore --host $MONGO_SERVER --drop --db $dir $dir"
	mongorestore --host $MONGO_SERVER --drop --db $dir $dir
done
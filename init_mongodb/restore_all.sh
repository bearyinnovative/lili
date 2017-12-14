#!/bin/bash

cd mongo_backup

dirs=(*)

for dir in "${dirs[@]}"
do
	echo "mongorestore --host $MONGO_SERVER --drop --db $dir $dir"
    mongorestore --host $MONGO_SERVER --drop --db $dir $dir
done
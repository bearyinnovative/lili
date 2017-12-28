#!/bin/bash


# 30 days ago
DATE=$(date -d '30 days ago' -u +"%Y-%m-%dT%H:%M:%SZ")
echo "Cleanup items haven't update since $DATE"

mongo lili --host $MONGO_SERVER --eval 'db.items.remove({updated: {$lt: ISODate("'"$DATE"'")}})'
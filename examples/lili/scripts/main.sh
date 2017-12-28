#!/bin/bash

case "$1" in
	backup) 
		./mongo-backup.sh
		;;

	cleanup) 
		./mongo-cleanup.sh
		;;

	restore) 
		./mongo-restore.sh
		;;

    *)
        echo $"choose one: {backup|cleanup|restore}"
        exit 1
esac
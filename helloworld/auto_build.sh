#!/bin/bash

cur=`pwd`
echo "please not close this process! monitor all source for auto build..."
while [ 1 = 1 ]
do
    build=0
    for i in `find $cur/ -name "*.go"`
    do
        i=${i//\/\//\/}
        backup="/tmp/$i"
        filename=`basename $i`
        dirname=`dirname $i`

        appdir="$cur/framework/app"
        appdir=${appdir//\/\//\/}
        if [ $dirname == $appdir ]
        then
            if [ $filename != "init.go" ] 
            then
                continue
            fi
        fi

        if [ -f "$backup" ] 
        then
            diff $backup $i > /tmp/t
            if [ $? -ne 0 ]
            then
                cp $i $backup
                echo "change:$i"
                build=1
                break
            fi
        else
            backup_dir=`dirname $backup`
            if [ ! -d $backup_dir ]
            then
                mkdir -p $backup_dir
            fi

            cp $i $backup
        fi
    done

    if [ $build -eq 1 ]
    then
        $cur/build.sh
        echo 
        echo "monitor all source for auto build..."
    fi
    sleep 0.1
done

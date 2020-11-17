#!/bin/bash
cd `dirname $0`
PROGRAM=supervise
PROGRAM_PATH=`pwd`


start() 
{ 
    stop; 
    cd $PROGRAM_PATH 
    echo -e "start $PROGRAM\c" 
    count=0 
    while [ $count -lt 5 ] 
    do 
        count=$((count+1)) 
        nohup $PROGRAM_PATH/supervise/$PROGRAM $PROGRAM_PATH 1>./run.log 2>&1 & 
        sleep 1 
        proc_count=`ps aux|grep $PROGRAM|grep $PROGRAM_PATH|grep -v grep|wc -l` 
        if [ $proc_count -eq 1 ] 
        then 
            echo  
            echo -e "$PROGRAM start \033[32m[ok]\033[0m" 
            return 
        else 
            echo -e ".\c" 
        fi 
    done 
    echo -e "$PROGRAM start \033[31m[failed]\033[0m" 
    cat ./run.log 
} 
 
stop() 
{ 
    echo -e "stop $PROGRAM\c" 
    count=0 
    while [ $count -lt 5 ] 
    do 
        count=$((count+1)) 
        proc_count=`ps aux|grep $PROGRAM|grep $PROGRAM_PATH|grep -v grep|wc -l` 
        if [ $proc_count -ge 1 ] 
        then 
            `ps aux|grep $PROGRAM|grep $PROGRAM_PATH|grep -v grep|awk '{print $2}'|xargs kill -9` 
            echo -e ".\c" 
            sleep 1 
        else  
            echo 
            echo -e "$PROGRAM stoped \033[32m[ok]\033[0m" 
            sh run stop
            return 
        fi 
    done 
    echo -e "$PROGRAM stoped \033[31m[failed]\033[0m" 
    exit -1 
} 
 
reopen() 
{ 
    echo -e "start reopen \c" 
    sh run reopen
} 

case "$1" in 
    'start')  
        start;; 
    'stop')  
        stop;; 
    'restart') 
        start;; 
    'reopen')
        reopen;;
    *) 
        echo "$0 (start|stop|restart)" 
esac 


#!/bin/bash

os=`uname -s`
if [ $os != "Linux" ]
    then
    export workspace=`pwd`
else
    fileLocation=$(readlink -f "$0")
    export workspace=$(dirname $fileLocation)
fi

version=v3
export version
export koala=$workspace"/.koala"
export program=HelloService

buildDir=$koala/"template/script/build"

function update_koala()
{
    supervise=$workspace"/supervise"

    if [ ! -d $koala ]
    then
        git clone https://micode.be.xiaomi.com/soa/koala.git
        mv koala $koala
    fi

    if [ ! -d $supervise ]
    then
          mkdir -p $supervise
    fi

    cd $koala
    git checkout master
    lastestOld=`git tag|grep -E "^$version"|sort -r|head -n 1`
    git pull
    lastest=`git tag|grep -E "^$version"|sort -r|head -n 1`
    if [ -z "$lastest" ]
    then
          echo "not found right tag of koala to build"
          exit -1
    fi
    if [ ! -z $lastestOld ] && [ $lastestOld != $lastest ]
    then
        if [ -d $koala/pkg ]
        then
            rm -rf $koala/pkg
        fi
    fi
    echo "use koala $lastest"
    git checkout $lastest -q
}

function main()
{
    update_koala
    source $buildDir"/run_plugin.sh"
    load_plugin_path

    run_plugin "${@}"
}

main "${@}"


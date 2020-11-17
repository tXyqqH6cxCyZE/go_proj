#!/bin/bash

repo=`git remote -v|tail -1 |awk '{print $2}'`
curPath=`pwd`
program=`basename $repo|awk -F '.' '{print $1}'`
packPath=~/.soa_pack
testBranch=SOA_TEST
productBranch=SOA_RELEASE
maxMinVersion=15
version=1.0
buildType=
binaryPath=
md5=
extra=()

function echo_failed()
{
    if [ $# -ne 1 ]
    then
        return
    fi

    echo -e "$1"
}

if [ $# -lt 1 ]
then
    echo "./package.sh all|test|product [customer directory...]"
    echo "for example:"
    echo "eg1 ./package.sh all ./data/"
    echo "eg2 ./package.sh all"
    exit
fi

if [ $1 != "test" ] && [ $1 != "all" ] && [ $1 != "product" ]
then
    echo "./package.sh all|test|product"
    exit
fi

untrackedFilesNum=`git status --porcelain 2>/dev/null | grep '^??' | grep -v '^??[[:space:]]\+bin/' | wc -l`
if [ $untrackedFilesNum -gt 0 ];then
    echo_failed "You have untracked files, please do a git-commit with following files or add these files to your .gitignore first:"
    git status --porcelain 2>/dev/null | grep '^??' | grep -v '^??[[:space:]]\+bin/' | sed 's/^??[[:space:]]*//g'
    exit
fi

diffStat=`git diff --name-only 2> /dev/null | grep -v '^bin/' | wc -l`
if [ $diffStat -gt 0 ];then
    echo_failed "Current branch is dirty, please do a git-commit with following files first:"
    git diff --name-only 2> /dev/null | grep -v '^bin/'
  exit
fi

currentRevision=$(git rev-parse HEAD)

buildType=$1

echo "current  repo:"$repo
shift
if [ $# -ge 1 ]
then
count=0
for i in "$@"
do
    extra[$((count))]=$i
    count=$((count+1))
done
fi

if [ -f "./build.sh" ]
then
    chmod +x ./build.sh
    ./build.sh linux

    if [ $? -ne 0 ]
    then
        echo "build failed..."
        exit
    fi
fi

function package_init()
{
    mkdir -p $packPath
}

function calc_md5()
{
    os=`uname -s`
    if [ $os != "Linux" ]
    then
        md5=`md5 $binaryPath|awk -F '=' '{print $2}'`
    else
        md5=`md5sum $binaryPath|awk '{print $1}'`
    fi
}

function create_version()
{
     if [ $# -ne 1 ]
    then
        echo_failed "invalid parameter, create_version path"
        exit
    fi

    build=`date '+%Y-%m-%d %H:%M:%S'`
    versionPath=$1/version
    if [ ! -f $versionPath ]
    then
        echo "version=$version" >> $versionPath
        echo "build=$build" >> $versionPath
        echo "md5=$md5" >> $versionPath
        echo "commit_id=$currentRevision" >> $versionPath
        return
    fi

    min=0
    max=1
    versionLine=
    while read LINE
    do
       count=`echo $LINE|grep version|wc -l`
       if [ $count -eq 0 ]
       then
           continue
       fi

       max=`echo $LINE|grep version|awk -F '=' '{print $2}'|awk -F '.' '{print $1}'`
       min=`echo $LINE|grep version|awk -F '=' '{print $2}'|awk -F '.' '{print $2}'`
    done < $versionPath

    min=$((min+1))
    if [ $min -ge $maxMinVersion ]
    then
        min=$((min%$maxMinVersion))
        max=$((max+1))
    fi

    if [  -z $max ]
    then
        max=1
    fi

    version="$max.$min"
    echo "version=$version" > $versionPath
    echo "build=$build" >> $versionPath
    echo "md5=$md5" >> $versionPath
}

function echo_succ()
{
    if [ $# -ne 1 ]
    then
        return
    fi

    echo -e "$1"
}

function check_group_and_service()
{
    if [ $# -ne 1 ]
    then
        echo_failed "invalid parameter in check_group_and_service"
        exit
    fi

    confPath=$1
    if [ ! -f $confPath ]
    then
        echo_failed "$confPath is not existed"
        exit
    fi

    group=
    service=
    while read LINE
    do
        count=`echo $LINE|grep group|wc -l`
       if [ $count -eq 1 ]
       then
           group=`echo $LINE|grep group|awk -F '=' '{print $2}'|awk -F '#' '{print $1}'`
           continue
       fi

       count=`echo $LINE|grep service|wc -l`
       if [ $count -eq 1 ]
       then
           service=`echo $LINE|grep service|awk -F '=' '{print $2}'|awk -F '#' '{print $1}'`
           continue
       fi
    done < $confPath

    if [ -z $group ] || [ -z $service ]
    then
        echo_failed "致命错误，请在$confPath中填写group[组名]和service[服务名]"
        exit
    fi
}

function remove_gitignore() {
    rm -rf $1/.gitignore
}

function package()
{
    if [ $# -ne 2 ]
    then
        echo_failed "invalid parameter in function package"
        exit
    fi

    cd $2
    mkdir -p $2/bin
    mkdir -p $2/config

    if [ ! -d $curPath/bin/ ]
    then
        echo_failed "致命错误，$curPath/bin不存在，无法打包"
        exit
    fi
    rm -v $2/bin/* 2>/dev/null

    count=`ls  $curPath/bin/*-linux 2>/dev/null |grep -v pid|wc -l`
    if [ $count -gt 1 ]
    then
        echo_failed "致命错误，在$curPath/bin/目录下存在多个二进制文件"
        exit
    elif [ $count -eq 1 ]
    then
        filename=`ls $curPath/bin/*-linux 2>/dev/null`
        filename=`basename $filename`
        destname=`echo $filename | awk -F '-linux' '{print $1}'`
        cp $curPath/bin/$filename $2/bin/$destname
    else
        count=`ls  $curPath/bin/|grep -v pid|wc -l`
        if [ $count -ne 1 ]
        then
            echo_failed "致命错误，在$curPath/bin/目录下存在多个二进制文件或不存在二进制文件"
            exit
        fi
        cp $curPath/bin/* $2/bin/
    fi

    if [ -f $2/bin/pid ]
    then
        rm $2/bin/pid
    fi
    binaryPath=$2/bin/`ls $2/bin/|grep -v pid`
    echo $binaryPath

    if [ -d $2/supervise ]
    then
        rm -rf $2/supervise
    fi
    mkdir -p $2/supervise
    cp $curPath/run $2/
    cp $curPath/run.sh $2/
    cp  $curPath/supervise/supervise $2/supervise

    for data in ${extra[@]}
    do
        echo_succ "正在打包:$data"
        filename=`basename $data`
        pathname=`dirname $data`
        mkdir -p $2/$pathname
        cp -rf "$curPath/$pathname/$filename" $2/$pathname
    done

    confType=test
    if [ $1 = "test" ]
    then
        cp $curPath/config/scm_config.ini.test.env $2/config/scm_config.ini
    else
        confType=online
        cp $curPath/config/scm_config.ini.online.env $2/config/scm_config.ini
    fi

    check_group_and_service "$curPath/config/scm_config.ini.$confType.env"

    calc_md5
    create_version $2
}

function package_program()
{
    if [ $# -ne 1 ]
    then
        echo_failed "invalid parameter in function package_program()"
        exit

    fi
    echo "package "$1
    fullPath=$packPath/$1/
    mkdir -p $fullPath

    if [ ! -d $curPath/.git ]
    then
        echo_failed "致命错误，$curPath不是git项目，请先在micode上创建项目"
        exit
    fi

    cd $fullPath
    if [ ! -d $fullPath/$program/.git ]
    then
        mkdir -p $fullPath/$program
        cd $fullPath/$program
        echo_succ "begin copy $repo"
        cp -rf $curPath/.git .
        echo_succ "copy from $repo succ"
    fi

    cd $fullPath/$program
    git checkout .
    git checkout master
    git pull
    branch=$testBranch
    if [ $1 = "product" ]
    then
        branch=$productBranch
    fi

    curBranch=`git branch|grep "\*"|awk '{print $2}'`

    if [ $curBranch != $branch ]
    then
        git checkout $branch -q
        if [ $? -ne 0 ]
        then
            echo_succ "create branch $branch"
            git checkout -b $branch -q
            if [ $? -ne 0 ]
            then
                echo_failed "checkout to $branch failed"
                exit
            fi

            rm -rf ./*
            git commit -a -m "first commit"
            git push origin $branch
            echo_succ "create branch $branch succ"
        fi
        git checkout .
    fi

    git pull origin $branch -q
    package $1 $fullPath/$program
    remove_gitignore $fullPath/$program
    git add -A --ignore-errors .
    git commit -a -m "package "$1 -q

    git push origin $branch -q

    echo_succ "package $1 [ok]"
    echo_succ "version $version"
    echo_succ "md5 $md5"
}

function package_all()
{
    package_init

    if [ $1 = "all" ]
    then
        package_program test
        package_program product
        return
    fi

    package_program $1
}

package_all $buildType

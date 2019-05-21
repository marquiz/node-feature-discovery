#!/bin/bash -e

this=`basename $0`
_ps1="\033[32;1muser@host\033[0m $ "
_print="pv -q"

usage()
{
    cat << EOF
Usage: $this [-h] [-d] SCRIPT_FILE

Options:
    -h      show this help
    -d      dry run
EOF
}

# Parse command line options
while getopts "hd" opt; do
    case $opt in
        h)  usage
            exit 0
            ;;
        d)  dry_run=1
            ;;
        *)  usage
            exit 1
            ;;
    esac
done

# Parse positional options
shift "$((OPTIND - 1))"
if [ $# -ne 1 ]; then
    usage
    exit 1
fi

cmd()
{
    echo -en "$_ps1"
    echo -n "$1" | $_print -L 20
    sleep 1
    echo
    if [ -z "$dry_run" ]; then
        sh -c "$1"
    else
        echo "..."
    fi
}

comment()
{
    echo -en "$_ps1 \\"
    echo -e "\033[31;1m"
    echo -n '>' "$1" | $_print -L 30
    sleep 1
    echo -e "\033[0m"
}


run_demo()
{
    cat $1 | while read line; do
        if [ "${line:0:1}" == "#" ]; then
            comment "$line"
        else
            cmd "$line"
        fi
    done

}

run_demo $1

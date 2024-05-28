#!/bin/bash -e

set -o pipefail

if [ $# -ne 4 ]; then
    echo -e "Usage: `basename $0` FROM TO COMMAND TEMPLATE\n"
    exit 1
fi

FROM=$1
TO=$2
CMD=$3
TEMPLATE=$4

do_cmd() {
    echo Processing from $1 to $2, cmd $3, template $4


    local yaml=`mktemp -p .`
    for i in `seq -f '%04.0f' $1 $2`; do
        sed s"/NODE_NAME/dummy-$i/" $4 >> "$yaml"
        echo --- >> "$yaml"
    done
    kubectl $3 -f "$yaml"
    rm "$yaml"
}


COUNT=$(($TO - $FROM))

BATCH_SIZE=10
BATCH_LAST=$(($COUNT / $BATCH_SIZE))

echo "Running $(($BATCH_LAST + 1)) batches of size $BATCH_SIZE"
for j in `seq 0 $BATCH_LAST`; do
    to=$(($FROM + $BATCH_SIZE - 1))
    [ "$to" -gt "$TO" ] && to=$TO
    do_cmd $FROM $to $CMD $TEMPLATE

    FROM=$(($FROM + $BATCH_SIZE))
done

#!/bin/bash
this_dir=`dirname $0`

export DEMO_PROMPT=1

ns="nfd"
if [ -n "$1" ]; then
    ns="$1"
fi

# Window watching NFD resources
#tmux new-session -d watch kubectl -n $ns get ds,svc,po -o wide
tmux new-session -d watch "$this_dir/watch-ns.sh" $ns

# Window watching node labels and annotations
# NOTE: just fire up bash to run watch-nodes.sh manually to work arounc some
#       tmux/watch/color issues
#tmux split-window -h -d watch --color "kubectl get no -o json | jq -C -r '.items[] | {name:.metadata.name, labels:.metadata.labels, annotations:.metadata.annotations}' | grep -v -e '\"beta.kubernetes.io' -e '.alpha.kubernetes.io' -e 'volumes.kubernetes.io' -e '\"kubernetes.io'"
tmux split-window -h -d bash

# "Main" interactive window
tmux split-window -v -d bash

# Window for watching demo resources
#tmux split-window -v -d watch kubectl -n demo get ds,svc,po -o wide
tmux split-window -v -d watch "$this_dir/watch-ns.sh" demo

tmux attach

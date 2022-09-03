#!/bin/bash

ERR_CODE_FILE=errcode.log
> $ERR_CODE_FILE # empty the file

N=1
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        N=$(grep -c 'cpu cores' /proc/cpuinfo | uniq)
elif [[ "$OSTYPE" == "darwin"* ]]; then
        N=$(sysctl -n hw.ncpu)
fi

if (($N <= 0)); then
    N=1
fi

echo "running ${N} concurrent plugin compilation."

mkdir -p dist/plugins

docompile() {
    go build -gcflags="all=-N -l" -buildmode=plugin -o ./dist/plugins/grpc-gateway-"$1".so ./"$2"
    echo $? >> $3
}

for d in ./internal/plugin/*/ ; do
    folder="$(basename $d)"
    echo ">> Building Plugin: $folder"
    docompile "$folder" "$d" "$ERR_CODE_FILE" &

    # allow to execute up to $N jobs in parallel
    if [[ $(jobs -r -p | wc -l) -ge $N ]]; then
        # now there are $N jobs already running, so wait here for any job
        # to be finished so there is a place to start next one.
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
          wait -n
        fi
    fi
done

# no more jobs to be started but wait for pending jobs
# (all need to be finished)
wait

# check for any compilation failures and report it accordingly
# so the Pipeline can react properly.

COUNT_FAILURES=$(grep -wv 0 -i $ERR_CODE_FILE | wc -l)

if (($COUNT_FAILURES > 0)); then
    echo "plugin compilation failed, there are $COUNT_FAILURES failure(s). Please kindly check the complete log."
    exit 1
fi
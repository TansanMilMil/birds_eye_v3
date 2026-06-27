#!/bin/bash -eu

cd `dirname $0`

if [ $BIRDSEYEAPI_EXECUTION_MODE = 'PRODUCTION' ] ; then
    echo 'running PRODUCTION mode...'
    ./go/dist/birds_eye_v3
else 
    echo 'running DEBUG mode...'
    # bashを対話モードで動かしてdockerが終了しないようにする
    /bin/bash
fi

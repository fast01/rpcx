#!/bin/bash
# Add 'syntax="proto2";' to the beginning of all proto files under 
# PWD(recursively) if the proto file does not have it.
for file in $(find . -name "*.proto"); do
    if grep -q 'syntax\s*=\s*"proto2";' $file; then
        echo "[already had] $file  , delete !";

        if  echo $(uname) |grep -q -i "darwin" ; then
            sed -i ".bak" '/syntax.*=.*;/d' $file
        else
            sed -i".bak"  '/syntax.*=.*;/d' $file
        fi
    else
        if  echo $(uname) |grep -q -i "darwin" ; then
            sed -i ".bak" '1s/^/syntax="proto2";'$'\\\n/' $file
        else
            sed -i".bak"  '1s/^/syntax="proto2";\n/' $file
        fi

        echo "[replaced] $file"
    fi
done

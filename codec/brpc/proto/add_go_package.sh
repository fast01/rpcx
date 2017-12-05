#!/bin/bash
# Add 'go_package="$*";' to the beginning of all proto files under
# PWD(recursively) if the proto file does not have it.
#

gp="$1"
gp=${gp//\//\\/}
#echo "$gp"
#exit 0
for file in $(find . -name "*.proto"); do
    if grep -q 'option\s*go_package\s*=\s*.*' $file; then
        echo "[already had] $file";
    else
        if  echo $(uname) |grep -q -i "darwin" ; then
            sed -i "" -E '1s/^/go_package='\"$gp\"$';\\\n/' $file
        else
            sed -i  '1s/^/go_package='\"$gp\"';\n/' $file
        fi
        echo "[replaced] $file"

    fi
done

#option go_package = "github.com/myuser/myprotos/person";
#!/bin/sh

set -e

examples_dir="./examples"
bin_dir="./bin"
target=$1

if [ -z $target ]
then
    target="examples"
fi

mkdir -p "$bin_dir"

if [ "$target" = "examples" ]
then 
    dir="$examples_dir/*"
    for sowo_file in $dir
    do
        name=$(basename $sowo_file '.sowo')
        c_file="$bin_dir/$name.c"
        go run . $sowo_file -o $c_file --save-tokens --save-ast
        cc -Wall $c_file -o "$bin_dir/$name"
    done
elif [ "$target" = "clean" ]
then
    dir="$examples_dir/*"
    for sowo_file in $dir
    do
        name=$(basename $sowo_file '.sowo')
        rm -f "$bin_dir/$name".c
        rm -f "$bin_dir/$name"_ast.json
        rm -f "$bin_dir/$name"_tok.txt
        rm -f "$bin_dir/$name"
    done
else
    echo "Unknown build target $target"
    exit 1
fi
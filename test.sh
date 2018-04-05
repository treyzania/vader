#!/bin/bash

if [ ! -f vader ]; then
	go build
fi

vaderbin=$(realpath vader)

for f in examples/*; do
	pushd $f > /dev/null
	echo "running test $f..."
	$vaderbin run
	echo -e
	popd > /dev/null
done

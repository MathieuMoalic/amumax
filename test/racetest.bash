#! /bin/bash

# builds with -race and runs tests with browser open.

set -e

go install -race github.com/MathieuMoalic/amumaxcmd/mumax3

google-chrome http://localhost:35367 &

for f in *.mx3; do
	mumax3 $f 
done

go install github.com/MathieuMoalic/amumaxcmd/mumax3 # re-build without race detector


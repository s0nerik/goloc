rm -rf .tmp
gox -output=".tmp/{{.OS}}_{{.Arch}}" -os="darwin linux windows" -arch="amd64"
cd .tmp
zip -r goloc.zip .
cd ..
mkdir out
mv .tmp/goloc.zip out/
rm -rf .tmp
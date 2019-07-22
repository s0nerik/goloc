rm -rf .tmp
gox -output=".tmp/{{.OS}}_{{.Arch}}" -os="darwin linux windows" -arch="amd64"

cd .tmp

cat <<EOT >> goloc.sh
#!/bin/bash

case "\$OSTYPE" in
  darwin*)  EXECUTABLE="darwin_amd64" ;;
  linux*)   EXECUTABLE="linux_amd64" ;;
  msys*)    EXECUTABLE="windows_amd64.exe" ;;
  *)
	  echo "Platform is not supported: \$OSTYPE"
	  exit 1
  ;;
esac

\${0%/*}/\${EXECUTABLE} "\$@"
EOT

zip -r goloc.zip .

cd ..
mkdir out
mv .tmp/goloc.zip out/
rm -rf .tmp
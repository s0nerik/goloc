rm -rf .tmp
gox -output=".tmp/{{.OS}}_{{.Arch}}" -os="darwin linux windows" -arch="amd64"

cd .tmp

cat <<EOT >> goloc.sh
#!/bin/bash

if [ -x "\$(command -v goloc)" ]; then
  echo "Using goloc from PATH..."
  echo "Version:"
  goloc --version
  goloc "\$@"
  exit 0
else
  case "\$OSTYPE" in
    darwin*)  EXECUTABLE="darwin_amd64" ;;
    linux*)   EXECUTABLE="linux_amd64" ;;
    msys*)    EXECUTABLE="windows_amd64.exe" ;;
    *)
      echo "Platform is not supported: \$OSTYPE"
      exit 1
    ;;
  esac

  echo "Using goloc from project..."
  echo "Version:"
  "\${0%/*}"/\${EXECUTABLE} --version
  "\${0%/*}"/\${EXECUTABLE} "\$@"
fi
EOT

zip -r goloc.zip .

cd ..
mkdir out
mv .tmp/goloc.zip out/
rm -rf .tmp
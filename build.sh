rm -rf .tmp
gox -output=".tmp/{{.OS}}_{{.Arch}}" -os="darwin linux windows" -arch="amd64 arm64"

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
  EXECUTABLE=""
  case "\$OSTYPE" in
    darwin*)  EXECUTABLE+="darwin_" ;;
    linux*)   EXECUTABLE+="linux_" ;;
    msys*)    EXECUTABLE+="windows_" ;;
    *)
      echo "Platform is not supported: \$OSTYPE"
      exit 1
    ;;
  esac
  
  MACHINE_TYPE=\$(uname -m)
  case "\$MACHINE_TYPE" in
    arm64*)    EXECUTABLE+="arm64" ;;
    x86_64*)   EXECUTABLE+="amd64" ;;
    *)
      echo "CPU Architecture is not supported: \$MACHINE_TYPE"
      exit 1
    ;;
  esac
  
  case "\$OSTYPE" in
    msys*) EXECUTABLE+=".exe"";;
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
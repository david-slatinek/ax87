#!/bin/bash

cert="<certificate path>"

ip=("<servers ip>")

clean() {
  rm -rf files/
  rm files.tar.bz2
}

init() {
  cd server-cert || exit

  mkdir files
  cp server-cert.pem files/
  cp server-key.pem files/
  cp ../../cert/ca-cert/ca-cert.pem files

  tar -cvjf files.tar.bz2 files/

  for x in "${ip[@]}"; do
    echo -e "\n\nUploading to $x\n"
    scp -i "$cert" files.tar.bz2 "ec2-user@$x":/home/ec2-user/
    echo -e "\n\n"

    ssh -i "$cert" -tt "ec2-user@$x" <<EOF
        sudo yum update -y

        rm -rf cert/
        mkdir -p cert/ca-cert

        rm -rf api/
        mkdir -p api/server-cert

        mv files.tar.bz2 api/
        cd api/
        tar -xvjf files.tar.bz2
        rm files.tar.bz2

        mv files/* .
        rm -rf files/

        mv ca-cert.pem ../cert/ca-cert

        mv *.pem server-cert/

        logout
EOF
  done

  clean
}

upload_api() {
  make compile

  mkdir -p files/env

  mv api files/
  cp env/.env files/env/

  tar -cvjf files.tar.bz2 files/

  for x in "${ip[@]}"; do
    echo -e "\n\nUploading to $x\n"
    scp -i "$cert" files.tar.bz2 "ec2-user@$x":/home/ec2-user/api/
    echo -e "\n\n"

    ssh -i "$cert" -tt "ec2-user@$x" <<EOF
      cd api || exit
      rm api
      rm -rf env/

      tar -xvjf files.tar.bz2
      rm files.tar.bz2
      mv files/* .
      rm -r files/

      ./api &

      logout
EOF
  done

  clean
}

help() {
  echo "Supported arguments:"
  echo "'init' - perform servers initialization"
  echo "'api' - upload api to servers"
  echo "'all' - perform 'init' and 'api'"
  echo -e "\nExample: ./upload.sh init"
}

if [ $# -eq 0 ]; then
  echo "No argument!"
  help
  exit 1
fi

if [[ $1 == "init" ]]; then
  init
elif [[ $1 == "api" ]]; then
  upload_api
elif [[ $1 == "all" ]]; then
  init
  upload_api
else
  echo "Invalid argument: $1!"
  help
  exit 1
fi

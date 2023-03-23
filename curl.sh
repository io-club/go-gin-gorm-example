#! /usr/bin/env bash

base_url="http://101.42.21.155:80/api"

case $1 in
1010) # register
  curl -X POST "$base_url"/auth/register -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
1011) # login
  curl -X POST "$base_url"/auth/login -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
2010) # get fabrics list
  curl -X GET "$base_url"/list
  ;;
2011) # get fabric
  curl -X GET $base_url/fabric/"$2"
  ;;
2012) # post fabric
  curl -X POST $base_url/fabric \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "image=@/home/trdthg/b.jpg"
  ;;
2013) # put fabric
  curl -X PUT $base_url/fabric/"$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@/home/trdthg/b.jpg"
  ;;
2014) # delete fabric
  curl -X DELETE $base_url/fabric/"$2"
  ;;
*)
  echo "Usage: $0 {get|post|put|delete}"
  ;;
esac

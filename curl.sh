#! /usr/bin/env bash

dev=0
if [ $dev -eq 1 ]; then
  base_url="http://localhost:8080/api"
else
  base_url="http://101.42.21.155:80/api"
fi

# 00       00
# service, api

case $1 in
1000) # register
  curl -X POST "$base_url"/auth/register -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
1001) # login
  curl -X POST "$base_url"/auth/login -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
2000) # get fabrics list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  #     - category: string optional, 例如: defult, new, hot, ...
  curl -X GET "$base_url/fabric/list?page=1&size=10&category=default"
  ;;
2001) # get fabric
  curl -X GET $base_url/fabric/"$2"
  ;;
2002) # post fabric
  curl -X POST $base_url/fabric \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "category=default" \
    -F "image=@/home/trdthg/resources/a.jpg"
  ;;
2003) # put fabric
  curl -X PUT $base_url/fabric/"$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@/home/trdthg/resources/a.jpg"
  ;;
2004) # delete fabric
  curl -X DELETE $base_url/fabric/"$2"
  ;;
3001) # upload images
  curl -X POST $base_url/image/upload \
    -F "tableName=fabrics" \
    -F "recordId=3" \
    -F "images=@/home/trdthg/resources/a.jpg" \
    -F "images=@/home/trdthg/resources/a.jpg" \
    -F "images=@/home/trdthg/resources/a.jpg"
  ;;
3002) # delete image
  curl -X DELETE $base_url/image/"$2"
  ;;
*)
  echo "Usage: $0 {get|post|put|delete}"
  ;;
esac

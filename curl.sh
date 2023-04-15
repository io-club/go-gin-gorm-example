#! /usr/bin/env bash

set -eu

BASE_URL="http://localhost:8080/api"
echo "🤣 MODE: $MODE"
if [ "${MODE}" = "release" ]; then
  BASE_URL="http://101.42.21.155:80/api"
fi
echo "🤣 base url: $BASE_URL"

set -eux

# 00       00
# service, api

case $1 in
1000) # register
  curl -X POST "$BASE_URL"/auth/register -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
1001) # login
  curl -X POST "$BASE_URL"/auth/login -d '{
    "username": "test1",
    "password": "123"
  }'
  ;;
2000) # get fabrics list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  #     - category: string optional, 例如: defult, new, hot, ...
  curl -X GET "$BASE_URL/fabric/list?page=1&size=10&category=default"
  ;;
2001) # get fabric
  curl -X GET "$BASE_URL/fabric/$2"
  ;;
2002) # post fabric
  curl -X POST $BASE_URL/fabric \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "category=default" \
    -F "image=@${2}"
  ;;
2003) # update fabric
  curl -X PUT "$BASE_URL/fabric/$2" \
    -F "name=aaa1"
  # -F "detail=bbb2" \
  # -F "image=@$3"
  ;;
2004) # delete fabric
  curl -X DELETE "$BASE_URL/fabric/$2"
  ;;
3001) # upload images
  # 上传图片
  # tableName: string required, 例如: fabric, brand, ...
  #     用于指定图片属于哪个表，参考各个 api 的命名规则
  # recordId: string required, 例如: 1，2, ...
  # images: file required, 例如: a.jpg, b.jpg, ...
  #     图片可以传多个，但是数据库最多只能传 5 张

  curl -X POST $BASE_URL/image/upload \
    -F "tableName=$2" \
    -F "recordId=$3" \
    -F "images=@$4"
  # -F "images=@/home/trdthg/resources/a.jpg" \
  # -F "images=@/home/trdthg/resources/a.jpg" \
  ;;
3002) # delete image
  curl -X DELETE "$BASE_URL/image/$2"
  ;;
4000) # get brand list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  curl -X GET "$BASE_URL/brand/list?page=1&size=10"
  ;;
4001) # get brand
  curl -X GET "$BASE_URL/brand/$2"
  ;;
4002) # create brand
  curl -X POST $BASE_URL/brand \
    -F "name=例如" \
    -F "detail=bbb" \
    -F "image=@${2}"
  ;;
4003) # update brand
  curl -X PUT "$BASE_URL/brand/$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@$3"
  ;;
4004) # delete brand
  curl -X DELETE "$BASE_URL/brand/$2"
  ;;
5000) # get trend list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  curl -X GET "$BASE_URL/trend/list?page=1&size=10"
  ;;
5001) # get trend
  curl -X GET "$BASE_URL/trend/$2"
  ;;
5002) # create trend
  curl -X POST $BASE_URL/trend \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "image=@$2"
  ;;
5003) # update trend
  curl -X PUT "$BASE_URL/trend/$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@$3"
  ;;
5004) # delete trend
  curl -X DELETE "$BASE_URL/trend/$2"
  ;;
6000) # get cloth list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  curl -X GET "$BASE_URL/cloth/list?page=1&size=10"
  ;;
6001) # get cloth
  curl -X GET "$BASE_URL/cloth/$2"
  ;;
6002) # create cloth
  curl -X POST $BASE_URL/cloth \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "image=@$2"
  ;;
6003) # update cloth
  curl -X PUT "$BASE_URL/cloth/$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@$3"
  ;;
6004) # delete cloth
  curl -X DELETE "$BASE_URL/cloth/$2"
  ;;
7000) # get dress list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  curl -X GET "$BASE_URL/dress/list?page=1&size=10"
  ;;
7001) # get dress
  curl -X GET "$BASE_URL/dress/$2"
  ;;
7002) # create dress
  curl -X POST $BASE_URL/dress \
    -F "name=aaa" \
    -F "detail=bbb" \
    -F "image=@$2"
  ;;
7003) # update dress
  curl -X PUT "$BASE_URL/dress/$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@$3"
  ;;
7004) # delete dress
  curl -X DELETE "$BASE_URL/dress/$2"
  ;;
8000) # get news list
  # query:
  #     - page: int required, 例如: 1, 2, 3, ...
  #     - size: int required, 例如: 10, 20, 30, ...
  curl -X GET "$BASE_URL/news/list?page=1&size=10"
  ;;
8001) # get news
  curl -X GET "$BASE_URL/news/$2"
  ;;
8002) # create news
  # type: industry | school_company
  curl -X POST $BASE_URL/news \
    -F "main=aaa" \
    -F "title=bbb" \
    -F "type=industry" \
    -F "image=@$2"
  ;;
8003) # update news
  curl -X PUT "$BASE_URL/news/$2" \
    -F "name=aaa" \
    -F "detail=bbb2" \
    -F "image=@$3"
  ;;
8004) # delete news
  curl -X DELETE "$BASE_URL/news/$2"
  ;;
*)
  echo "Usage: $0 {get|post|put|delete}"
  ;;
esac

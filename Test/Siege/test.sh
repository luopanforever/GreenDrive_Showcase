#!/bin/bash

# 设置测试的URL
URL="http://localhost:8080/car/show/car1/scene.bin"
echo "start testing"
# 并发用户数组
CONCURRENCIES=(1 5 10 50 100 200 500)

# 测试每个并发级别
for CONCURRENCY in "${CONCURRENCIES[@]}"
do
  echo "Testing with $CONCURRENCY users..."
  siege -c $CONCURRENCY -r 10 --log=siege.log $URL
done

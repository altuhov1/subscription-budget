#!/bin/bash

URL="http://localhost:8080"
USER_ID="550e8400-e29b-41d4-a716-446655440001"

echo "📍 Тестируем $URL"
echo "👤 User ID: $USER_ID"
echo

echo "1. /api/subs/create (Первая подписка)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/subs/create" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Netflix","price":100,"user_id":"'"$USER_ID"'","start_date":"01-2025"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
SUB_ID_1=$(echo "$body" | jq -r '.numberOfSub' 2>/dev/null)
echo

echo "2. /api/subs/create (Вторая подписка)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/subs/create" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Spotify","price":200,"user_id":"'"$USER_ID"'","start_date":"02-2025"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
SUB_ID_2=$(echo "$body" | jq -r '.numberOfSub' 2>/dev/null)
echo

echo "3. /api/get-subs/param (Подсчет двух подписок)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/param" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'","start_date":"01-2025","subscriptions":["Netflix","Spotify"]}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

echo "4. /api/get-subs/user-id/one (Первая подписка)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/one" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID_1"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

echo "5.  /api/get-subs/user-id/one (Вторая подписка)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/one" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID_2"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

echo "6. /api/get-subs/user-id/all (Ожидаем 2 подписки)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/all" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

echo "7. /api/subs/update (Обновляем Netflix)"
resp=$(curl -s -w "\n%{http_code}" -X PUT "$URL/api/subs/update" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID_1"',"service_name":"Netflix Premium","price":150,"user_id":"'"$USER_ID"'","start_date":"03-2025"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
[ -n "$body" ] && echo "   → Ответ: $body" || echo "   → Тело: (пусто, как и ожидается для 204)"
echo

echo "8. /api/get-subs/param (Подсчет после обновления)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/param" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'","start_date":"01-2025","subscriptions":["Netflix Premium","Spotify"]}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

echo "9. /api/subs/delete (Удаляем Netflix Premium)"
resp=$(curl -s -w "\n%{http_code}" -X DELETE "$URL/api/subs/delete" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID_1"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
[ -n "$body" ] && echo "   → Ответ: $body" || echo "   → Тело: (пусто)"
echo

echo "10. /api/subs/delete (Удаляем Spotify)"
resp=$(curl -s -w "\n%{http_code}" -X DELETE "$URL/api/subs/delete" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID_2"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
[ -n "$body" ] && echo "   → Ответ: $body" || echo "   → Тело: (пусто)"
echo

echo "11. /api/get-subs/user-id/all (Проверка что подписок не осталось)"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/all" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

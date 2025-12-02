#!/bin/bash

URL="http://localhost:8080"
USER_ID="550e8400-e29b-41d4-a716-446655440002"

echo Тестируем $URL"
echo "User ID: $USER_ID"
echo

# 1. Создание
echo "1. 📥 /api/subs/create"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/subs/create" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Test","price":100,"user_id":"'"$USER_ID"'","start_date":"01-2025"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
SUB_ID=$(echo "$body" | jq -r '.numberOfSub' 2>/dev/null)
echo

# 2. Расчёт
echo "2. 🧮 /api/get-subs/param"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/param" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'","start_date":"02-2025","end_date":"12-2025","subscriptions":["Test"]}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

# 3. Получение по ID
echo "3. 🔍 /api/get-subs/user-id/one"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/one" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

# 4. Список по user_id
echo "4. 📋 /api/get-subs/user-id/all"
resp=$(curl -s -w "\n%{http_code}" -X POST "$URL/api/get-subs/user-id/all" \
  -H "Content-Type: application/json" \
  -d '{"user_id":"'"$USER_ID"'"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
echo "   → Ответ: $body"
echo

# 5. Обновление
echo "5. ✏️ /api/subs/update"
resp=$(curl -s -w "\n%{http_code}" -X PUT "$URL/api/subs/update" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID"',"service_name":"Updated","price":150,"user_id":"'"$USER_ID"'","start_date":"02-2025"}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
[ -n "$body" ] && echo "   → Ответ: $body" || echo "   → Тело: (пусто, как и ожидается для 204)"
echo

# 6. Удаление
echo "6. 🗑️ /api/subs/delete"
resp=$(curl -s -w "\n%{http_code}" -X DELETE "$URL/api/subs/delete" \
  -H "Content-Type: application/json" \
  -d '{"id":'"$SUB_ID"'}')
http_code=$(echo "$resp" | tail -n1)
body=$(echo "$resp" | sed '$d')
echo "   → Код: $http_code"
[ -n "$body" ] && echo "   → Ответ: $body" || echo "   → Тело: (пусто)"
echo

echo "✅ Готово. Проверьте вывод выше."
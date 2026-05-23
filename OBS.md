
# error add item

2026/05/21 17:49:51 /app/internal/infrastructure/database/shopping_cart_db.go:70 ERROR: insert or update on table "shopping_cart_items" violates foreign key constraint "fk_shopping_carts_items" (SQLSTATE 23503)
[1.757ms] [rows:0] UPDATE "shopping_carts" SET "created_at"='2026-05-21 17:49:51.069',"updated_at"='2026-05-21 17:49:51.069',"user_id"='34bcdf2d-e5f0-4bf5-9e1e-f5e53c018e64',"total_btc"=0.01 WHERE "id" = '10e4e4cb-dedf-4b21-be4c-16cd5f34aaac'
2026/05/21 17:49:51 [51034b094749/EjOHHcG3Uy-000006] "POST http://localhost:8080/api/v1/cart/items HTTP/1.1" from 172.18.0.1:51918 - 500 66B in 11.501919ms
[ERROR] 2026/05/21 17:49:51 logger.go:184: [51034b094749/EjOHHcG3Uy-000006] internal error (500): ERROR: insert or update on table "shopping_cart_items" violates foreign key constraint "fk_shopping_carts_items" (SQLSTATE 23503)
2026/05/21 17:50:04 [51034b094749/EjOHHcG3Uy-000007] "GET http://app:8080/metrics HTTP/1.1" from 172.18.0.5:43032 - 200 2623B in 1.867431ms

# error delete cart

2026/05/21 17:51:19 /app/internal/infrastructure/database/shopping_cart_db.go:27 record not found
[0.188ms] [rows:0] SELECT * FROM "shopping_carts" WHERE user_id = '34bcdf2d-e5f0-4bf5-9e1e-f5e53c018e64' ORDER BY "shopping_carts"."id" LIMIT 1
2026/05/21 17:51:19 [51034b094749/EjOHHcG3Uy-000014] "DELETE http://localhost:8080/api/v1/cart/ HTTP/1.1" from 172.18.0.1:55098 - 500 66B in 1.027607ms
[ERROR] 2026/05/21 17:51:19 logger.go:184: [51034b094749/EjOHHcG3Uy-000014] internal error (500): record not found


# error update item remove item, car not found 

2026/05/21 17:52:07 /app/internal/infrastructure/database/shopping_cart_db.go:27 record not found
[0.202ms] [rows:0] SELECT * FROM "shopping_carts" WHERE user_id = '34bcdf2d-e5f0-4bf5-9e1e-f5e53c018e64' ORDER BY "shopping_carts"."id" LIMIT 1
2026/05/21 17:52:07 [51034b094749/EjOHHcG3Uy-000020] "PATCH http://localhost:8080/api/v1/cart/items/9079e549-5615-41e3-9115-c9ed39616835 HTTP/1.1" from 172.18.0.1:55098 - 404 37B in 842.064µs

2026/05/21 17:52:11 /app/internal/infrastructure/database/shopping_cart_db.go:27 record not found
[0.238ms] [rows:0] SELECT * FROM "shopping_carts" WHERE user_id = '34bcdf2d-e5f0-4bf5-9e1e-f5e53c018e64' ORDER BY "shopping_carts"."id" LIMIT 1
2026/05/21 17:52:11 [51034b094749/EjOHHcG3Uy-000021] "DELETE http://localhost:8080/api/v1/cart/items/9079e549-5615-41e3-9115-c9ed39616835 HTTP/1.1" from 172.18.0.1:55098 - 404 37B in 953.331µs
^C
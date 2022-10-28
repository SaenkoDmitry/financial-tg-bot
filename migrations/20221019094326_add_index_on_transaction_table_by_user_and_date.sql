-- +goose Up
-- +goose StatementBegin
CREATE INDEX transaction_user_id_created_at_idx ON route256.financial_bot.transaction USING BTREE (user_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX route256.financial_bot.transaction_user_id_created_at_idx;
-- +goose StatementEnd


-- КАК ВЫБРАЛ?
-- посмотрел выдачу EXPLAIN ANALYZE по запросу:
-- EXPLAIN (ANALYZE) SELECT category_id, amount, created_at FROM financial_bot.transaction
-- 			WHERE user_id = 123548568 AND created_at > '2022-03-05 12:04:05.588718';
-- понял, что тут будет много операций сравнения по timestamp, поэтому
-- обычное дерево можно использовать и оно должно дать выигрыш по перфомансу
-- применение индекса судя по ответу планировщика postgres дало это ускорение
-- и в explain analyze также видно, что планировщик использует индекс в запросе

-- ниже выкладки до и после добавления индекса:

---- BEFORE APPLYING INDEX
-- Seq Scan on transaction  (cost=0.00..2310.68 rows=11258 width=22) (actual time=0.020..26.445 rows=11344 loops=1)
--   Filter: ((created_at > '2022-03-05 12:04:05.588718'::timestamp without time zone) AND (user_id = 123548568))
--   Rows Removed by Filter: 88701
-- Planning Time: 0.089 ms
-- Execution Time: 27.111 ms

-- AFTER APPLYING INDEX
-- Bitmap Heap Scan on transaction  (cost=159.69..1138.56 rows=11258 width=22) (actual time=0.902..7.477 rows=11344 loops=1)
--   Recheck Cond: ((user_id = 123548568) AND (created_at > '2022-03-05 12:04:05.588718'::timestamp without time zone))
--   Heap Blocks: exact=810
--   ->  Bitmap Index Scan on created_at_timestamp_idx  (cost=0.00..156.87 rows=11258 width=0) (actual time=0.768..0.768 rows=11344 loops=1)
--         Index Cond: ((user_id = 123548568) AND (created_at > '2022-03-05 12:04:05.588718'::timestamp without time zone))
-- Planning Time: 0.111 ms
-- Execution Time: 7.991 ms
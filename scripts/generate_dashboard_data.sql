-- 生成仪表盘测试数据

-- 清理旧的测试订单数据
DELETE FROM inbound_items WHERE created_at >= '2025-01-01';
DELETE FROM outbound_items WHERE created_at >= '2025-01-01';
DELETE FROM inbound_orders WHERE created_at >= '2025-01-01';
DELETE FROM outbound_orders WHERE created_at >= '2025-01-01';

-- 插入入库订单（最近30天，每天多笔）
INSERT INTO inbound_orders (order_no, supplier_id, warehouse_id, status, total_quantity, created_at, created_by, updated_at, updated_by)
SELECT 
    CONCAT('RK', DATE_FORMAT(date_val, '%Y%m%d'), LPAD(FLOOR(RAND() * 1000), 3, '0')),
    FLOOR(1 + RAND() * 9),  -- supplier_id: 1-9
    FLOOR(1 + RAND() * 7),  -- warehouse_id: 1-7
    FLOOR(1 + RAND() * 2),  -- status: 1-2
    FLOOR(50 + RAND() * 200),  -- total_quantity: 50-250
    date_val,
    1,
    date_val,
    1
FROM (
    SELECT DATE_ADD('2026-04-09', INTERVAL seq DAY) as date_val
    FROM (
        SELECT a.N + b.N * 10 + c.N * 100 as seq
        FROM 
            (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) a
            ,(SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) b
            ,(SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2) c
    ) numbers
    WHERE seq < 30
) dates
WHERE RAND() > 0.3;  -- 70%的天数有订单

-- 插入出库订单（最近30天，每天多笔）
INSERT INTO outbound_orders (order_no, customer_id, warehouse_id, status, total_quantity, created_at, created_by, updated_at, updated_by)
SELECT 
    CONCAT('CK', DATE_FORMAT(date_val, '%Y%m%d'), LPAD(FLOOR(RAND() * 1000), 3, '0')),
    FLOOR(1 + RAND() * 3),  -- customer_id: 1-3
    FLOOR(1 + RAND() * 7),  -- warehouse_id: 1-7
    FLOOR(1 + RAND() * 2),  -- status: 1-2
    FLOOR(30 + RAND() * 150),  -- total_quantity: 30-180
    date_val,
    1,
    date_val,
    1
FROM (
    SELECT DATE_ADD('2026-04-09', INTERVAL seq DAY) as date_val
    FROM (
        SELECT a.N + b.N * 10 + c.N * 100 as seq
        FROM 
            (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) a
            ,(SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) b
            ,(SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2) c
    ) numbers
    WHERE seq < 30
) dates
WHERE RAND() > 0.3;  -- 70%的天数有订单

-- 为每个入库订单添加明细
INSERT INTO inbound_items (inbound_order_id, product_id, location_id, quantity, batch_no, created_at, created_by, updated_at, updated_by)
SELECT 
    io.id,
    FLOOR(1 + RAND() * 16),  -- product_id: 1-16
    FLOOR(1 + RAND() * 21),  -- location_id: 1-21
    FLOOR(10 + RAND() * 50),  -- quantity: 10-60
    CONCAT('BATCH', DATE_FORMAT(io.created_at, '%Y%m%d'), FLOOR(RAND() * 100)),
    io.created_at,
    1,
    io.created_at,
    1
FROM inbound_orders io
WHERE io.created_at >= '2025-01-01';

-- 为每个出库订单添加明细（每个订单3-5个产品）
INSERT INTO outbound_items (outbound_order_id, product_id, location_id, quantity, created_at, created_by, updated_at, updated_by)
SELECT 
    oo.id,
    FLOOR(1 + RAND() * 16),  -- product_id: 1-16
    FLOOR(1 + RAND() * 21),  -- location_id: 1-21
    FLOOR(5 + RAND() * 30),  -- quantity: 5-35
    oo.created_at,
    1,
    oo.created_at,
    1
FROM outbound_orders oo
WHERE oo.created_at >= '2025-01-01'
HAVING RAND() > 0.6;  -- 平均每个订单1-2个产品

-- 再添加一些产品到出库订单（使订单更真实）
INSERT INTO outbound_items (outbound_order_id, product_id, location_id, quantity, created_at, created_by, updated_at, updated_by)
SELECT 
    oo.id,
    FLOOR(1 + RAND() * 16),
    FLOOR(1 + RAND() * 21),
    FLOOR(5 + RAND() * 30),
    oo.created_at,
    1,
    oo.created_at,
    1
FROM outbound_orders oo
WHERE oo.created_at >= '2025-01-01'
HAVING RAND() > 0.5;

-- 插入今天和昨天的订单（确保今天有数据）
INSERT INTO inbound_orders (order_no, supplier_id, warehouse_id, status, total_quantity, created_at, created_by, updated_at, updated_by)
VALUES 
    ('RK20260509001', 3, 1, 2, 150, NOW(), 1, NOW(), 1),
    ('RK20260509002', 5, 2, 2, 80, NOW(), 1, NOW(), 1),
    ('RK20260509003', 7, 1, 1, 200, NOW(), 1, NOW(), 1);

INSERT INTO outbound_orders (order_no, customer_id, warehouse_id, status, total_quantity, created_at, created_by, updated_at, updated_by)
VALUES 
    ('CK20260509001', 1, 1, 2, 100, NOW(), 1, NOW(), 1),
    ('CK20260509002', 2, 2, 2, 60, NOW(), 1, NOW(), 1),
    ('CK20260509003', 3, 1, 1, 80, NOW(), 1, NOW(), 1);

-- 为今天的入库订单添加明细
INSERT INTO inbound_items (inbound_order_id, product_id, location_id, quantity, batch_no, created_at, created_by, updated_at, updated_by)
VALUES 
    ((SELECT id FROM inbound_orders WHERE order_no = 'RK20260509001'), 6, 1, 50, 'BATCH2026050901', NOW(), 1, NOW(), 1),
    ((SELECT id FROM inbound_orders WHERE order_no = 'RK20260509001'), 7, 2, 100, 'BATCH2026050902', NOW(), 1, NOW(), 1),
    ((SELECT id FROM inbound_orders WHERE order_no = 'RK20260509002'), 16, 5, 80, 'BATCH2026050903', NOW(), 1, NOW(), 1),
    ((SELECT id FROM inbound_orders WHERE order_no = 'RK20260509003'), 11, 3, 200, 'BATCH2026050904', NOW(), 1, NOW(), 1);

-- 为今天的出库订单添加明细
INSERT INTO outbound_items (outbound_order_id, product_id, location_id, quantity, created_at, created_by, updated_at, updated_by)
VALUES 
    ((SELECT id FROM outbound_orders WHERE order_no = 'CK20260509001'), 6, 1, 30, NOW(), 1, NOW(), 1),
    ((SELECT id FROM outbound_orders WHERE order_no = 'CK20260509001'), 8, 2, 70, NOW(), 1, NOW(), 1),
    ((SELECT id FROM outbound_orders WHERE order_no = 'CK20260509002'), 16, 5, 60, NOW(), 1, NOW(), 1),
    ((SELECT id FROM outbound_orders WHERE order_no = 'CK20260509003'), 12, 4, 80, NOW(), 1, NOW(), 1);

-- 统计插入的数据
SELECT '入库订单' as type, COUNT(*) as count FROM inbound_orders WHERE created_at >= '2025-01-01'
UNION ALL
SELECT '出库订单', COUNT(*) FROM outbound_orders WHERE created_at >= '2025-01-01'
UNION ALL
SELECT '入库明细', COUNT(*) FROM inbound_items WHERE created_at >= '2025-01-01'
UNION ALL
SELECT '出库明细', COUNT(*) FROM outbound_items WHERE created_at >= '2025-01-01';

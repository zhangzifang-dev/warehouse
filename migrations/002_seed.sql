-- Seed data for warehouse management system

-- Default admin user (password: admin123, hashed with bcrypt)
INSERT INTO users (id, username, password_hash, status, created_at, created_by, updated_at, updated_by) VALUES
(1, 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.rsQ5pPjZ5yVlWK5WAe', 1, NOW(), 1, NOW(), 1);

-- Default roles
INSERT INTO roles (id, name, description, status, created_at, created_by, updated_at, updated_by) VALUES
(1, '超级管理员', '拥有所有权限', 1, NOW(), 1, NOW(), 1),
(2, '仓库管理员', '管理仓库、货位、库存', 1, NOW(), 1, NOW(), 1),
(3, '采购员', '管理供应商、入库单', 1, NOW(), 1, NOW(), 1),
(4, '销售员', '管理客户、出库单', 1, NOW(), 1, NOW(), 1),
(5, '库管员', '库存查询、盘点', 1, NOW(), 1, NOW(), 1),
(6, '审计员', '查看审计日志', 1, NOW(), 1, NOW(), 1);

-- Assign super_admin role to admin user
INSERT INTO user_roles (user_id, role_id, created_at, created_by, updated_at, updated_by) VALUES
(1, 1, NOW(), 1, NOW(), 1);

-- Default permissions
INSERT INTO permissions (id, name, code, resource, action, created_at, created_by, updated_at, updated_by) VALUES
-- User management
(1, '查看用户列表', 'user:list', 'user', 'list', NOW(), 1, NOW(), 1),
(2, '创建用户', 'user:create', 'user', 'create', NOW(), 1, NOW(), 1),
(3, '编辑用户', 'user:update', 'user', 'update', NOW(), 1, NOW(), 1),
(4, '删除用户', 'user:delete', 'user', 'delete', NOW(), 1, NOW(), 1),
-- Role management
(5, '查看角色列表', 'role:list', 'role', 'list', NOW(), 1, NOW(), 1),
(6, '创建角色', 'role:create', 'role', 'create', NOW(), 1, NOW(), 1),
(7, '编辑角色', 'role:update', 'role', 'update', NOW(), 1, NOW(), 1),
(8, '删除角色', 'role:delete', 'role', 'delete', NOW(), 1, NOW(), 1),
-- Warehouse management
(9, '查看仓库列表', 'warehouse:list', 'warehouse', 'list', NOW(), 1, NOW(), 1),
(10, '创建仓库', 'warehouse:create', 'warehouse', 'create', NOW(), 1, NOW(), 1),
(11, '编辑仓库', 'warehouse:update', 'warehouse', 'update', NOW(), 1, NOW(), 1),
(12, '删除仓库', 'warehouse:delete', 'warehouse', 'delete', NOW(), 1, NOW(), 1),
-- Location management
(13, '查看货位列表', 'location:list', 'location', 'list', NOW(), 1, NOW(), 1),
(14, '创建货位', 'location:create', 'location', 'create', NOW(), 1, NOW(), 1),
(15, '编辑货位', 'location:update', 'location', 'update', NOW(), 1, NOW(), 1),
(16, '删除货位', 'location:delete', 'location', 'delete', NOW(), 1, NOW(), 1),
-- Category management
(17, '查看分类列表', 'category:list', 'category', 'list', NOW(), 1, NOW(), 1),
(18, '创建分类', 'category:create', 'category', 'create', NOW(), 1, NOW(), 1),
(19, '编辑分类', 'category:update', 'category', 'update', NOW(), 1, NOW(), 1),
(20, '删除分类', 'category:delete', 'category', 'delete', NOW(), 1, NOW(), 1),
-- Product management
(21, '查看商品列表', 'product:list', 'product', 'list', NOW(), 1, NOW(), 1),
(22, '创建商品', 'product:create', 'product', 'create', NOW(), 1, NOW(), 1),
(23, '编辑商品', 'product:update', 'product', 'update', NOW(), 1, NOW(), 1),
(24, '删除商品', 'product:delete', 'product', 'delete', NOW(), 1, NOW(), 1),
-- Inventory management
(25, '查看库存', 'inventory:list', 'inventory', 'list', NOW(), 1, NOW(), 1),
(26, '库存调整', 'inventory:adjust', 'inventory', 'adjust', NOW(), 1, NOW(), 1),
(27, '库存盘点', 'inventory:check', 'inventory', 'check', NOW(), 1, NOW(), 1),
-- Supplier management
(28, '查看供应商列表', 'supplier:list', 'supplier', 'list', NOW(), 1, NOW(), 1),
(29, '创建供应商', 'supplier:create', 'supplier', 'create', NOW(), 1, NOW(), 1),
(30, '编辑供应商', 'supplier:update', 'supplier', 'update', NOW(), 1, NOW(), 1),
(31, '删除供应商', 'supplier:delete', 'supplier', 'delete', NOW(), 1, NOW(), 1),
-- Customer management
(32, '查看客户列表', 'customer:list', 'customer', 'list', NOW(), 1, NOW(), 1),
(33, '创建客户', 'customer:create', 'customer', 'create', NOW(), 1, NOW(), 1),
(34, '编辑客户', 'customer:update', 'customer', 'update', NOW(), 1, NOW(), 1),
(35, '删除客户', 'customer:delete', 'customer', 'delete', NOW(), 1, NOW(), 1),
-- Inbound order
(36, '查看到货单列表', 'inbound:list', 'inbound', 'list', NOW(), 1, NOW(), 1),
(37, '创建到货单', 'inbound:create', 'inbound', 'create', NOW(), 1, NOW(), 1),
(38, '编辑到货单', 'inbound:update', 'inbound', 'update', NOW(), 1, NOW(), 1),
(39, '删除到货单', 'inbound:delete', 'inbound', 'delete', NOW(), 1, NOW(), 1),
(40, '确认入库', 'inbound:confirm', 'inbound', 'confirm', NOW(), 1, NOW(), 1),
-- Outbound order
(41, '查看出库单列表', 'outbound:list', 'outbound', 'list', NOW(), 1, NOW(), 1),
(42, '创建出库单', 'outbound:create', 'outbound', 'create', NOW(), 1, NOW(), 1),
(43, '编辑出库单', 'outbound:update', 'outbound', 'update', NOW(), 1, NOW(), 1),
(44, '删除出库单', 'outbound:delete', 'outbound', 'delete', NOW(), 1, NOW(), 1),
(45, '确认出库', 'outbound:confirm', 'outbound', 'confirm', NOW(), 1, NOW(), 1),
-- Stock transfer
(46, '查看调拨单列表', 'transfer:list', 'transfer', 'list', NOW(), 1, NOW(), 1),
(47, '创建调拨单', 'transfer:create', 'transfer', 'create', NOW(), 1, NOW(), 1),
(48, '编辑调拨单', 'transfer:update', 'transfer', 'update', NOW(), 1, NOW(), 1),
(49, '删除调拨单', 'transfer:delete', 'transfer', 'delete', NOW(), 1, NOW(), 1),
(50, '确认调拨', 'transfer:confirm', 'transfer', 'confirm', NOW(), 1, NOW(), 1),
-- Audit log
(51, '查看审计日志', 'audit:list', 'audit', 'list', NOW(), 1, NOW(), 1);

-- Assign all permissions to super_admin role
INSERT INTO role_permissions (role_id, permission_id, created_at, created_by, updated_at, updated_by)
SELECT 1, id, NOW(), 1, NOW(), 1 FROM permissions;

-- Test warehouses
INSERT INTO warehouses (id, name, address, status, created_at, created_by, updated_at, updated_by) VALUES
(1, '主仓库', '上海市浦东新区张江高科技园区', 1, NOW(), 1, NOW(), 1),
(2, '分仓库', '上海市嘉定区工业园区', 1, NOW(), 1, NOW(), 1);

-- Test locations
INSERT INTO locations (id, warehouse_id, zone, shelf, level, position, status, created_at, created_by, updated_at, updated_by) VALUES
(1, 1, 'A', '01', '1', '01', 1, NOW(), 1, NOW(), 1),
(2, 1, 'A', '01', '1', '02', 1, NOW(), 1, NOW(), 1),
(3, 1, 'A', '01', '2', '01', 1, NOW(), 1, NOW(), 1),
(4, 1, 'B', '02', '1', '01', 1, NOW(), 1, NOW(), 1),
(5, 2, 'C', '01', '1', '01', 1, NOW(), 1, NOW(), 1);

-- Test categories
INSERT INTO categories (id, name, parent_id, created_at, created_by, updated_at, updated_by) VALUES
(1, '电子产品', 0, NOW(), 1, NOW(), 1),
(2, '办公用品', 0, NOW(), 1, NOW(), 1),
(3, '手机', 1, NOW(), 1, NOW(), 1),
(4, '电脑', 1, NOW(), 1, NOW(), 1),
(5, '文具', 2, NOW(), 1, NOW(), 1);

-- Test products
INSERT INTO products (id, sku, name, category_id, spec, unit, created_at, created_by, updated_at, updated_by) VALUES
(1, 'SKU001', 'iPhone 15 Pro', 3, '256GB 深空黑色', '台', NOW(), 1, NOW(), 1),
(2, 'SKU002', 'MacBook Pro 14', 4, 'M3 Pro 18GB 512GB', '台', NOW(), 1, NOW(), 1),
(3, 'SKU003', 'AirPods Pro 2', 3, '主动降噪', '副', NOW(), 1, NOW(), 1),
(4, 'SKU004', '签字笔', 5, '0.5mm 黑色', '盒', NOW(), 1, NOW(), 1),
(5, 'SKU005', 'A4打印纸', 5, '80g 500张', '包', NOW(), 1, NOW(), 1);

-- Test suppliers
INSERT INTO suppliers (id, name, contact, phone, address, created_at, created_by, updated_at, updated_by) VALUES
(1, '华东电子供应商', '王经理', '13900000001', '上海市黄浦区', NOW(), 1, NOW(), 1),
(2, '北方办公用品', '赵经理', '13900000002', '北京市朝阳区', NOW(), 1, NOW(), 1);

-- Test customers
INSERT INTO customers (id, name, contact, phone, address, created_at, created_by, updated_at, updated_by) VALUES
(1, '科技公司A', '陈总', '13800000001', '上海市徐汇区', NOW(), 1, NOW(), 1),
(2, '贸易公司B', '刘总', '13800000002', '上海市静安区', NOW(), 1, NOW(), 1);

-- Test inventory
INSERT INTO inventory (warehouse_id, product_id, location_id, quantity, batch_no, created_at, created_by, updated_at, updated_by) VALUES
(1, 1, 1, 100, 'BATCH2024001', NOW(), 1, NOW(), 1),
(1, 2, 2, 50, 'BATCH2024002', NOW(), 1, NOW(), 1),
(1, 3, 3, 200, 'BATCH2024003', NOW(), 1, NOW(), 1),
(1, 4, 4, 500, 'BATCH2024004', NOW(), 1, NOW(), 1),
(1, 5, 4, 300, 'BATCH2024005', NOW(), 1, NOW(), 1),
(2, 1, 5, 30, 'BATCH2024006', NOW(), 1, NOW(), 1),
(2, 3, 5, 80, 'BATCH2024007', NOW(), 1, NOW(), 1);

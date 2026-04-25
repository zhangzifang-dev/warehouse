-- Drop Audit Table
DROP TABLE IF EXISTS audit_logs;

-- Drop Order Tables
DROP TABLE IF EXISTS stock_transfers;
DROP TABLE IF EXISTS outbound_items;
DROP TABLE IF EXISTS outbound_orders;
DROP TABLE IF EXISTS inbound_items;
DROP TABLE IF EXISTS inbound_orders;

-- Drop Core Business Tables
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS warehouses;

-- Drop User Permission Tables
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;

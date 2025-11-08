-- bcrypt("password") のハッシュ値（生成済）
-- ハッシュ: "$2a$10$ZbUe7P8nABe7/5A2IXxGcuJr1oX7v7rT4D.DTlh.q6t8EN3F1bG4m"

INSERT INTO users (email, password_hash)
VALUES ('test@example.com', '$2a$10$Q4A86sZk6FTXTZTonDVMm.npTH5yp2e8/vmvk2EWKLOGxmaPF127a');
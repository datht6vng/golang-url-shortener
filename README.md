* Thuật toán rút gọn:
    - Hash: SHA256
    - Sinh chuỗi id: gouuid
    - Base save hash: Base58
    --> Lấy 8 bit đầu tiên
* Tạo user để debug server local:
CREATE USER `root`@`%` IDENTIFIED BY "123456"; GRANT ALL PRIVILEGES ON *.* TO `root`@`%` WITH GRANT OPTION; FLUSH PRIVILEGES;

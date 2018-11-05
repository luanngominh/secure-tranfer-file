# Secure Transfer File - Network Security Project 1
# Protocol
Sau khi client và server thực hiện tcp handshark thành công:
* Bước 1:  Server sẽ gửi public key cho client, message có dạng
Key: <Public key>\n
* Bước 2: Client sinh ra client key. Key này được sinh ra bằng cách, client sẽ random một chuổi bất kỳ, sau đó hash bằng MD5, chuổi sau khi hash sẽ là key. Key này sẽ được encrypt bằng public key vừa nhận được từ server ở bước 1, sau đó gửi cho server. Message có dạng
Key: <Client key>\n
* Bước 3: Server decrypt client key bằng private key ⇒ lấy được client key. Server sẽ sinh ra session key sau đó encrypt session key bằng client key - encrypt bằng AES - và gửi session key cho client. Messes có dạng
Session: <session key>\n
* Bước 4: Client nhận được session key ở dạng cipher text, client dùng client key để giải mã ⇒ nhận được session key. Sau đó client sẽ gửi file request, message này được encrypt bằng AES với key là client key. Message có dạng
File: <Tên file>\n
Session: <Session key vừa nhận được>\n
* Bước 5: Server nhận được file request, giải mã bằng client key nhận được file request. So sánh session key, nếu trùng thì qua bước 6, ngược lại đóng kết nối.
* Bước 6: Server sẽ kiểm tra file có tồn tại hay không, nếu không có thì đóng kết nối. Nếu có, server sẽ đọc file, sau đó encrypt bằng client key và gửi qua client
* Bước 7: Client nhân nội dung file, dùng client key để decrypt

# Usage
## Server
Docker version is avaiable at `docker pull luanngominh/secure-tranfer-file`
### Env
SERVER_PORT=<Server run on port>
ADDR=<Listen on interface>
FILE_STORAGE=<File folder>
PUBLIC=<public key with base64 encoding>
PRIVATE=<private key with base64 encoding>

Example:
SERVER_PORT="1212"
ADDR=<"127.0.0.1"
FILE_STORAGE="${PWD}/files"
PUBLIC="blah blah =="
PRIVATE="blah balh =="

## Client
In demo, client will connect to localhost:1212
./client <tên file cần nhận>

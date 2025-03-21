# syntax=docker/dockerfile:1

##
## STEP 1 - BUILD
##

# Chỉ định image cơ sở là Golang, sử dụng phiên bản alpine để có một image nhẹ
# Đây là bước đầu tiên để xây dựng môi trường ứng dụng
FROM golang:1.24.1-alpine as builder

# Tạo thư mục làm việc trong container
WORKDIR /app

# Sao chép tất cả các file trong thư mục hiện tại (ví dụ như các file .go) vào thư mục làm việc của container
COPY . .

# Cài đặt công cụ swag để tự động tạo tài liệu API
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Khởi tạo tài liệu API bằng swag, với file main.go làm điểm khởi đầu
RUN swag init -g cmd/api/main.go

# Biên dịch ứng dụng Go và tạo file thực thi (binary) cho API
RUN go build -o api ./cmd/api

##
## STEP 2 - DEPLOY
##

# Sử dụng image Alpine nhẹ để triển khai ứng dụng
FROM alpine:latest

# Cài đặt ca-certificates và tzdata để đảm bảo tính bảo mật và đúng múi giờ
RUN apk update
RUN apk add --no-cache ca-certificates tzdata && update-ca-certificates

# Thiết lập múi giờ cho container là Asia/Ho_Chi_Minh (giờ Việt Nam)
ENV TZ=Asia/Ho_Chi_Minh
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Đặt thư mục làm việc trong container
WORKDIR /app

# Sao chép file thực thi (binary) 'api' từ bước build vào thư mục làm việc trong container
COPY --from=builder /app/api .

# Mở cổng 80 để ứng dụng có thể lắng nghe các yêu cầu từ bên ngoài
EXPOSE 80

# Chạy ứng dụng API khi container được khởi động
ENTRYPOINT ["./api"]

##
## Những phần bắt buộc cần có trong Dockerfile:
##

# 1. `FROM`: Đây là lệnh bắt buộc trong Dockerfile. Bạn phải chỉ định một image cơ sở (base image) mà container của bạn sẽ xây dựng từ đó.
#    - Ví dụ: `FROM golang:1.24.1-alpine` hoặc `FROM node:14`.

# 2. `WORKDIR`: Đặt thư mục làm việc trong container. Đây là nơi mà tất cả các lệnh tiếp theo sẽ được thực thi.
#    - Nếu thư mục không tồn tại, Docker sẽ tự tạo.

# 3. `COPY` hoặc `ADD`: Sao chép các file từ hệ thống của bạn vào container.
#    - `COPY . .` sẽ sao chép tất cả các file trong thư mục hiện tại vào thư mục làm việc trong container.

# 4. `RUN`: Lệnh thực thi các câu lệnh trong container, ví dụ như cài đặt phần mềm hoặc biên dịch mã nguồn.
#    - `RUN go install github.com/swaggo/swag/cmd/swag@latest` cài đặt công cụ `swag`.

# 5. `ENTRYPOINT` hoặc `CMD`: Xác định lệnh mặc định khi container được khởi động. `ENTRYPOINT` thường dùng để xác định lệnh cố định mà container sẽ chạy.
#    - `ENTRYPOINT ["./api"]` chỉ định rằng khi container bắt đầu, ứng dụng `./api` sẽ được chạy.

# 6. `EXPOSE`: Mở cổng trong container cho phép các yêu cầu từ bên ngoài truy cập vào ứng dụng.
#    - `EXPOSE 80` mở cổng 80 để ứng dụng web có thể lắng nghe yêu cầu HTTP.


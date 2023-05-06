#                   BUILD STAGE

# mendefinisikan base image, karna project ini menggunakan golang maka kita memerlukan golang base image. 
# jika kita ingin file menjadi kecil kita gunakan image yang berakhiran alpine sama seperti image postgres
# FROM = specify base image
FROM golang:1.20.4-alpine3.17 AS builder

# WORKDIR: declare current working directory inside the image
WORKDIR /app

# COPY : copy necessery files
#  first . everything from current folder when we ru docker build command to  build an immage
# second . current working directory inside the image where file and folder are being copied to. so /app will be the place to copy the data
COPY . . 

# RUN : build our app to single binary executable file 
RUN go build -o main main.go


# setelah binary file diproduce kita akan lanjut ke stage selanjutnya
#                           RUN STAGE
FROM alpine:3.17
WORKDIR /app
# copy executable binary file from the BUILDER STAGE To this run stage image, the path we want to copy, target location the final image to copied that file to
COPY --from=builder /app/main .
# app.env ini hanya untuk dev saja tidak untu production. ini hanya sekedar contoh
COPY app.env . 

# EXPOSE : inform docker that the container listen on the specify network port at runtime. dalam kasus ini adalah 8080 sama seperti apa yang kita tulis pada file environtment
# note  instruksi expose aslinya tidak perish the port, ini hanya berfungsi sebagaii dokumentasi antara orang yang membuat image docker dengan orang yang akan menjalankan container tentang port mana yang akan digunakan untuk publish
EXPOSE 8080

# define the default command  to run when the container done
# dalam kasus ini kita hanya ingin menjalankan executable file yang kita buat di langkah sebelumnya jadi hanya ada ssatu nilai didalam kurung kurawal
CMD ["/app/main"]




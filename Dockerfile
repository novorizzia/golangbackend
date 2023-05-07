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

# install curl in builder state image, karna secara default base alpine image tidak memiliki preinstalled curl
RUN apk add curl

# melakukan migrate pada database
# hampir mirip dengan yang ada pada github ci.yml
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz


# setelah binary file diproduce kita akan lanjut ke stage selanjutnya
#                           RUN STAGE
FROM alpine:3.17
WORKDIR /app
# copy executable binary file from the BUILDER STAGE To this run stage image, the path we want to copy, target location the final image to copied that file to
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate

# app.env ini hanya untuk dev saja tidak untu production. ini hanya sekedar contoh
COPY app.env . 

COPY start.sh .
COPY wait-for.sh .

# copy semua file migration dari db/migration ke image
COPY db/migration ./migration

# EXPOSE : inform docker that the container listen on the specify network port at runtime. dalam kasus ini adalah 8080 sama seperti apa yang kita tulis pada file environtment
# note  instruksi expose aslinya tidak perish the port, ini hanya berfungsi sebagaii dokumentasi antara orang yang membuat image docker dengan orang yang akan menjalankan container tentang port mana yang akan digunakan untuk publish
EXPOSE 8080

# define the default command  to run when the container done
# dalam kasus ini kita hanya ingin menjalankan executable file yang kita buat di langkah sebelumnya jadi hanya ada ssatu nilai didalam kurung kurawal
CMD ["/app/main"]
# we will use entry point instruction
ENTRYPOINT [ "/app/start.sh" ]

# ketika CMD digunakan bersama ENTRYPOINT akan bertindak selayakanya additional parameter that will be pass to entrypoint script
# jadi seperti menjalankan ENTRYPOINT [ "/app/start.sh", "/app/main" ]
# tapi dengan pemisahan ini kita bisa lebih fleksible untuk mennggantinya dengan command lain saat runtime kapanpun kita mau 
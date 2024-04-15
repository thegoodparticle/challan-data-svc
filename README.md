docker build -t challan-data-svc:v1.0.0 .

docker-compose up

curl localhost:8081/challan-info/KA20AB1234
docker build -t challan-data-svc:v1.0.0 .

docker-compose up

curl localhost:8081/challan-info/KA20AB1234


$ challan-data-svc % kubectl create namespace challan
namespace/challan created

$ challan-data-svc % kubectl --context=minikube apply -f ./deployment/dynamodb.yaml -n=challan
persistentvolume/dynamodb-volume created
persistentvolumeclaim/dynamodb-volume-claim created
deployment.apps/dynamodb created
service/dynamodb created

$ challan-data-svc % kubectl get pods -n challan
NAME                        READY   STATUS    RESTARTS   AGE
dynamodb-554dc669bd-cds7n   1/1     Running   0          6m7s


eval  $(minikube docker-env)

docker build -t challan-data-svc:v1.1.0 .

$ challan-data-svc % kubectl --context=minikube apply -f ./deployment/challan-data-svc.yaml -n=challan 
secret/challan-data-svc-secret created
deployment.apps/challan-data-svc created
service/challan-data-svc created

$ challan-data-svc % kubectl get pods -n=challan
NAME                               READY   STATUS    RESTARTS   AGE
challan-data-svc-b9b86db44-vtb8q   1/1     Running   0          76s
dynamodb-554dc669bd-cds7n          1/1     Running   0          18m

{"Vehicle Number": "KA20EF9012","Unit Name": "Cyberabad","Date": "02-Apr-2024","Time": "15:33","Place of Violation": "RAIDURGAM TR PS LIMITS","PS Limits": "Raidurgam Tr PS","Violation": "Wrong Parking in the carriage way","Fine Amount": 100}

# Skyfarm-backend


Where Response is «all info about..» - it’s very long response, so i don’t insert it


## GetNodeInfo send the list of the pods in namespace :
**GET Request : localhost:3000/api/kube_get/:namespace**
Response : list of pods on node in given namespace (now i got only 1 pod «mongo» in default
namespace) :
>{
"node name": "minikube",
"pods names": [
"mongo"
]
}

## CreateNamespace creates a new namespace:
**POST Request : localhost:3000/api/kube/create_namespace**
Body : just string with name of namespace, for example : customm
Response : all information about created namespace (metadata, spec, status)
CreatePV creates new perststent volume:
POST Request : localhost:3000/api/kube/create_pv
Body :
>{
"name" : "mongo-pv",
"storage" : "1Gi",
"path" : "/tmp/mongodb"
}

Response : all info about created pv (metadata, spec, status)
## CreatePVC creates new perststent volume claim:
**POST Request : localhostL3000/api/kube/create_pvc**
Body :
>{
"name" : "mongo-pvg",
"storage" : "1Gi",
"namespace" : "custom"
}

Response : all info about created pvc (metadata, spec, status)
## CreateNodePort — creates new nodeport:
**POST Request : localhost:3000/api/kube/create_nodeport**
Body :
>{
"name" : "testport",
"namespace" : "default",
"port" : 27017,
"redirect_port" : 30001
}

Response : all info about created pvc (metadata, spec, status)
## CreateOrUpdateConfigmap — creates or updates data of configmap if it already exists:
**POST Request : localhost:3000/api/kube/create_config_map**
Body :
>{
"name": "conff",
"namespace" : "default",
"env" : {
"MESSAGE" : "Hello from Moeid",
"NEWMESSAGE" : "Godfather",
"PASSED" : "FALSE"
}
}

Response : «created/updated»: "conff"
## CreateOrUpdateSecret — creates or updates data of secret if it already exists:
**POST Request : localhost:3000/api/kube/create_secret**
Body :
>{
"name": "secrr",
"namespace": "custom",
"env": {
"MONGO_INITDB_ROOT_USERNAME" : "admin",
"MONGO_INITDB_ROOT_PASSWORD" : "admin"
}
}

Response : «created/updated»: "secrr"
## CreatePod — creates new pod:
**POST Request : localhost:3000/api/kube/add**
Body :
>{
"command" : [],
"configmap_name" : "conff",
"port" : 27017,
"secret_name" : "secrr",
"name" : "mongo2",
"namespace" : "default",
"container_name" : "mongo-container",
"claim_name" : "mongo-pvc",
"image" : "mongo",
"volume_name" : "mongodb-data",
"mountpath" : "/usr/share/mongo"
}

image — image of container. Could be image from docker hub
Response : all data about created pod (metadata, spec, containers, status)

## DeletePod — deletes pod from namespace:
**DELETE Request : localhost:3000/api/kube_delete/:namespace/:pod_name**
Response : «message» : «pod_name is deleted»

## CreateRole - creates new kubernetes role:
**POST Request localhost:3000/api/kube/create_role**

Body : 

>{
    "name":"role1",
    "namespace":"default", 
    "resourses" : ["pods"], 
    "verbs" : ["get", "list", "watch"]
}

response : all info about role 

## CreateServiceAccount - creates new service account 

**POST Request: localhost:3000/api/kube/create_account** 

Body : 
>{
    "name":"serv",
    "namespace":"default", 
    "secret-namespace" : "default", 
    "secret-name" : "secrr" 
}

response : all info about service account
## Create RoleBind - creates new role bind 

**POST Request localhost:3000/api/kube/role_bind**

Body : 
>{
    "name":"bind3",
    "role-name":"role1",
    "account-name":"serv",
    "namespace":"default"
}

response : all info about rolebind
## HelmGetCharts — get all information about installed charts:
**GET Request : localhost:3000/api/helm_get**
Response : all info about charts(name. Info. Status, metadata)
## HelmCreateChart — creates chart from url or from installed repository:
**POST Request: localhost:3000/api/helm**
Body if chart in repo :
>{
"chart_path":"bitnami/keycloak",
"namespace":"default",
"release_name":"good-keycloack"
}

Body if chart path is url :
>{
"chart_path":"https://charts.bitnami.com/bitnami/keycloak-13.0.2.tgz",
"namespace":"default",
"release_name":"good-keycloack"
}

## Helm Create Repo - creates helm chart repository by url:
**POST Request localhost:3000/api/helm_create_repo**
Body :
>{
"name" : "stable",
"url" : "https://charts.helm.sh/stable"
}

response : repo created
## Websocket connection — send the status of pod and events in default namespace :
**ws:localhost:12121/**
when you send message with name of pod it started to send this pod status every second and events
in default namespace when they’re happening/
# Database requests (already been in boilerplate, just changed to mongo):
## Create Test :
Creates new element in database
**POST Request localhost:3000/api/test**
Body :
>{
"Name": "Godfather",
"Passed": true,
"Number": 228
}

Response : Test Created
## Get Test :
Sends all elements in database
**GET Request localhost:3000/api/test**
Response : list of all elements
Get One Test :
Sends one element from database
GET Request localhost:3000/api/test/:id
Response :
>{
"ID": "63e38bee89f28de8cecb51c0",
"Name": "Godfather",
"Passed": true,
"Number": 228,
"CreatedAt": "2023-02-08T11:47:58.395Z",
"UpdatedAt": "2023-02-08T11:47:58.395Z"
}

## Update Test :
**POST Request localhost:3000/api/test/:id**
Updates element in database
Body :
>{
"Name": "GodFATHER",
"Passed": true,
"Number": 245
}

Response : Test updated
## Delete Test :
**DELETE Request localhost:3000/api/test/:id**
**Doesn’t work now, should delete elements from database by id. Elements doesn’t deleting**
Response : Test Deleted

## Get Auth code — made to simplify getting acces token from keycloak
don’t works whenyou’re not authorized in keycloak, and works only in browser, because of redirection
if follow this link (localhost:8080 is keycloack), keycloack will redirect to server, server will
make post request to keycloack using jwt secret key and code from keycloack, gets the token
and send it.
http://localhost:8080/realms/master/protocol/openid-connect/auth?client_id=skyfarm&response_type=code&redirect_uri=http://localhost:3000/get_code
Response : access token
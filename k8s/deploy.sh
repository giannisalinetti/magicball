#!/bin/sh

NAMESPACE=magicball
DB_NAME=appdb
APP_NAME=magicball
DB_DEPLOYMENT_CFG=appdb-deployment.yml
DB_SVC_CFG=appdb-svc.yml
APP_DEPLOYMENT_CFG=magicball-deployment.yml
APP_SVC_CFG=magicball-svc.yml

kubectl create namespace $NAMESPACE

echo "# Initializing database"
kubectl apply -f $DB_DEPLOYMENT_CFG -n $NAMESPACE
kubectl apply -f $DB_SVC_CFG -n $NAMESPACE

echo "# Waiting for database startup"
STATUS=$(kubectl get pods -n $NAMESPACE | grep $DB_NAME | awk '{print $2}')
while [ $STATUS != '1/1' ]; do
    sleep 5
    STATUS=$(kubectl get pods -n $NAMESPACE | grep $DB_NAME | awk '{print $2}')
done

echo "# Initializing application"
kubectl apply -f $APP_DEPLOYMENT_CFG -n $NAMESPACE
kubectl apply -f $APP_SVC_CFG -n $NAMESPACE

echo "# Waiting for application startup"
STATUS=$(kubectl get pods -n $NAMESPACE | grep $APP_NAME | awk '{print $2}')
while [ $STATUS != '1/1' ]; do
    sleep 5
    STATUS=$(kubectl get pods -n $NAMESPACE | grep $APP_NAME | awk '{print $2}')
done

exit 0

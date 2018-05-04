#!/bin/sh

NAMESPACE=magic8
APP_SVC=magicball

APP_URL=$(minikube service $APP_SVC -n $NAMESPACE --url)

printf "Application $APP_SVC is available at:\t$APP_URL\n"

exit 0

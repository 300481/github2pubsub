#!/usr/bin/env bash

ENVIRONMENT=(GITHUB_SECRET GITHUB_EVENTS GCP_TOPIC_NAME GCP_CREATE_TOPIC GCP_PROJECT_ID)

FUNCTION_PATH="$(dirname $(readlink -f $0))/../cloud-function"

die() {
    RC=$1
    shift
    echo ""
    echo "$@"
    exit ${RC}
}

for VARIABLE in ${ENVIRONMENT[@]} ; do
    echo -n "Checking variable '${VARIABLE}'..."
    [[ -z ${!VARIABLE} ]] && die 1 "Variable '${VARIABLE}' not set."
    echo "OK"
done

# login if not done already
CONFIG_CONTAINER_NAME=gcloud-config
if ! docker ps -a | grep -q ${CONFIG_CONTAINER_NAME} ; then
    docker run -ti --name ${CONFIG_CONTAINER_NAME} google/cloud-sdk:alpine gcloud auth login
fi

# set project id
docker run --rm -ti \
    --volumes-from ${CONFIG_CONTAINER_NAME} google/cloud-sdk:alpine \
        gcloud config set project ${GCP_PROJECT_ID}

# install cloud function
docker run --rm -ti \
    --volumes-from ${CONFIG_CONTAINER_NAME} \
    --volume ${FUNCTION_PATH}:/function \
    google/cloud-sdk:alpine \
        gcloud functions deploy ${GCP_TOPIC_NAME} \
        --runtime=go111 \
        --region=europe-west1 \
        --memory=128MB \
        --entry-point Send \
        --trigger-http \
        --set-env-vars=GCP_PROJECT_ID=${GCP_PROJECT_ID},GCP_TOPIC_NAME=${GCP_TOPIC_NAME},GCP_CREATE_TOPIC=${GCP_CREATE_TOPIC},GITHUB_SECRET=${GITHUB_SECRET},GITHUB_EVENTS=${GITHUB_EVENTS} \
        --source /function

echo ""
echo ""
echo ""

# get cloud function trigger URL
docker run --rm -ti \
    --volumes-from ${CONFIG_CONTAINER_NAME} google/cloud-sdk:alpine \
        gcloud functions describe ${GCP_TOPIC_NAME} \
        --region=europe-west1 \
        --format="value(httpsTrigger.url)"

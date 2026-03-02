#!/bin/bash

# Ensure configured gcloud locally and have created
# a GCP project
export GCP_PROJECT_ID=""

export SA_NAME="github-deploy-service"
# export GIT_REPOSITORY_OWNER="Hamfri"
export GIT_REPO="Hamfri/social"
export POOL_ID="social-github-deploy"
export PROVIDER_ID="github"
export ARTIFACT_REPO="social-repo"

gcloud config set project ${GCP_PROJECT_ID}

#########################################################################
# Create workload identity pool
#########################################################################

gcloud iam workload-identity-pools create ${POOL_ID} \
    --location="global" \
    --display-name="GitHub Actions Pool"


#########################################################################
# Create the OIDC Provider for GitHub
#########################################################################

gcloud iam workload-identity-pools providers create-oidc ${PROVIDER_ID} \
    --location="global" \
    --workload-identity-pool=${POOL_ID} \
    --display-name="GitHub Provider" \
    --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.actor=assertion.actor,attribute.aud=assertion.aud" \
    --attribute-condition="assertion.repository=='${GIT_REPO}'" \
    --issuer-uri="https://token.actions.githubusercontent.com"


#########################################################################
# Create the Service Account
#########################################################################

gcloud iam service-accounts create ${SA_NAME} \
    --display-name="GitHub Actions Service Account"

export SA_EMAIL=$(gcloud iam service-accounts list \
    --filter="displayName:GitHub Actions Service Account" \
    --format='value(email)')


#########################################################################
# Grant permissions and roles to service account
#########################################################################


PROJECT_NUMBER=$(gcloud projects describe ${GCP_PROJECT_ID} --format='value(projectNumber)')

# Impersonation roles 
gcloud iam service-accounts add-iam-policy-binding ${SA_EMAIL} \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/attribute.repository/${GIT_REPO}"


gcloud iam service-accounts add-iam-policy-binding ${SA_EMAIL} \
    --role="roles/iam.serviceAccountTokenCreator" \
    --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/attribute.repository/${GIT_REPO}"


# Project level deployment roles
roles=(
    "roles/artifactregistry.admin"
    "roles/artifactregistry.writer"
    "roles/cloudfunctions.admin"
    "roles/run.developer"
    "roles/iam.serviceAccountAdmin"
    "roles/iam.serviceAccountUser"
    "roles/iam.serviceAccountTokenCreator"
)

for role in "${roles[@]}"; do
    gcloud projects add-iam-policy-binding ${GCP_PROJECT_ID} \
    --member="serviceAccount:${SA_EMAIL}" \
    --role="${role}"
done

gcloud services enable iamcredentials.googleapis.com

#########################################################################
# Create Docker repository
#########################################################################

gcloud artifacts repositories create ${ARTIFACT_REPO} \
    --repository-format=docker \
    --location=europe-west1 \
    --description="Docker repository for social-api"


#########################################################################
# DEBUG 
#########################################################################

# Printout WIF_PROVIDER (COPY this to Github Secrets)
# gcloud iam workload-identity-pools providers list \
#   --workload-identity-pool="${POOL_ID}" \
#   --location=global \
#   --format="value(name)"

# WIF_SERVICE_ACCOUNT = ${SA_EMAIL} (Copy this to Github secrets too)

# check active roles
# gcloud iam service-accounts describe ${SA_EMAIL}
# gcloud iam service-accounts get-iam-policy ${SA_EMAIL} \
#     --format="json"


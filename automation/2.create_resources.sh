echo "First let's enable all APIs needed for this solution."
gcloud services enable storage.googleapis.com

echo "Now let's create a service account to access the resources with"
gcloud iam service-accounts create "$NAME-service" \
    --description="Service account to manage $NAME resources" \
    --display-name="$NAME Service"

echo "Now let's give the account the right role access to the project $PROJECT"
gcloud projects add-iam-policy-binding $PROJECT \
    --member="serviceAccount:$NAME-service@$PROJECT.iam.gserviceaccount.com" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT \
    --member="serviceAccount:$NAME-service@$PROJECT.iam.gserviceaccount.com" \
    --role="roles/run.admin"

echo "Creating storage bucket..."
gcloud alpha storage buckets create gs://$BUCKET_NAME --location $LOCATION

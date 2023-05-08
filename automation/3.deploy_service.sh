cd ../services/svont-svc

gcloud builds submit --tag eu.gcr.io/$PROJECT/$SERVICE_NAME

gcloud run deploy $SERVICE_NAME --image eu.gcr.io/$PROJECT/$SERVICE_NAME \
  --platform managed --project $PROJECT --region $REGION --allow-unauthenticated\
  --memory=256Mi --cpu=1 --service-account "$NAME-service@$PROJECT.iam.gserviceaccount.com" \
  --timeout 1800s --min-instances 1 --max-instances 1 --no-cpu-boost \
  --update-env-vars "BUCKET_NAME=$BUCKET_NAME,BUCKET_PATH=$BUCKET_PATH"

export CLOUD_RUN_URL=$(gcloud run services describe $SERVICE_NAME --platform managed --region $REGION --format 'value(status.url)')
# Save BUCKET_NAME for future inits

cd ../../automation

sed -i "/export CLOUD_RUN_URL=/c\export CLOUD_RUN_URL=$CLOUD_RUN_URL" 1.1.env_reinit.sh
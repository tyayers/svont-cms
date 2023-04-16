cd ../clients/test_vite_svelte

echo "Replace variables in the client code..."
sed -i "/VITE_CMS_SERVICE=/c\VITE_CMS_SERVICE=$CLOUD_RUN_URL" .env.production

echo "Building client..."
npm install
npm run build

echo "Now we're going to login to Firebase, please follow the instructions to login..."
firebase login --no-localhost

echo "Now deploying client to firebase..."
firebase projects:addfirebase $PROJECT || true

FIREBASE_NAME=${BUCKET_NAME//_/-}

echo "Setting Firebase config..."
#firebase init hosting
sed -i "/    \"target\":/c\    \"target\": \"$FIREBASE_NAME\"," firebase.json
sed -i "/    \"default\":/c\    \"default\": \"$PROJECT\"" .firebaserc

firebase use $PROJECT
firebase hosting:sites:create $FIREBASE_NAME
firebase target:apply hosting $FIREBASE_NAME $FIREBASE_NAME

echo "Deploying to Firebase..."
firebase deploy --only hosting:$FIREBASE_NAME
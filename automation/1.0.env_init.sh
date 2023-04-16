export PROJECT="cloud32x"
# Save PROJECT for future inits
sed -i "/export PROJECT=/c\export PROJECT=$PROJECT" 1.1.env_reinit.sh

export REGION="europe-west1"
# Save BUCKET_NAME for future inits
sed -i "/export REGION=/c\export REGION=$REGION" 1.1.env_reinit.sh

export NAME="cms"
# Save NAME for future inits
sed -i "/export NAME=/c\export NAME=$NAME" 1.1.env_reinit.sh

BUCKET_NAME="cms_$(tr -dc A-Za-z0-9 </dev/urandom | head -c 8 ; echo '')"
export BUCKET_NAME=$(echo "$BUCKET_NAME" | tr '[:upper:]' '[:lower:]')
# Save BUCKET_NAME for future inits
sed -i "/export BUCKET_NAME=/c\export BUCKET_NAME=$BUCKET_NAME" 1.1.env_reinit.sh
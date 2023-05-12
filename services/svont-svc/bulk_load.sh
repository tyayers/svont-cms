export TOKEN=

export HOST=https://cms2-service-ghfontasua-ew.a.run.app/posts
# export HOST=http://localhost:8080/posts

export CONTENT="A snowflake is a single crystal that has achieved a sufficient size, and may have amalgamated with others, which falls through the  as . Each flake nucleates around a dust particle in  air masses by attracting  cloud water droplets, which  and accrete in crystal"

for i in {1..10000}
do
   echo "Loop number $i"

   curl -X POST -i $HOST \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "title=test$i&content=$CONTENT&summary=$CONTENT&authorDisplayName=LoadTest&authorProfilePic=https://lh3.googleusercontent.com/a/AEdFTp7CoLZ5_FFs4MT_m5yxNBZI9FbTT8K2dOcCRdQNZQ=s96-c"
   
   sleep .1
done
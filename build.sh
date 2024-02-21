docker buildx build --platform linux/amd64 -t kawalrealcount:latest .
docker tag kawalrealcount alfianisnan26/kawalrealcount:latest
docker push alfianisnan26/kawalrealcount:latest

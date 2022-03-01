work_dir=$(dirname "$0")
cd ${work_dir}/..
docker-compose -f "docker-compose.yml" build
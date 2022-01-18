work_dir=$(dirname "$0")
cd ${work_dir}/..
docker-compose -f docker-compose.yml logs -f -t --tail="50"
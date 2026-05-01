#bin/bash
echo "Delete DB Data"
docker exec -it city-explorer-db psql -U cityexplorer -d city_explorer -c "DELETE FROM places WHERE source='overpass';"
echo "GET NYC FOOD"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 5
echo "completed trace script"
sleep 10
echo "GET NYC FOOD AGAIN"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 50
echo "GET NYC FOOD ONE MORE TIME"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 10
echo "GET WRONG URL LINKS"
curl http://localhost:8080/cities/new%20york%20city/cafe
curl http://localhost:8080/cities/new%20york%20city/build
curl http://localhost:8080/cities/new%20york%20city/
sleep 100



echo "Delete DB Data"
docker exec -it city-explorer-db psql -U cityexplorer -d city_explorer -c "DELETE FROM places WHERE source='overpass';"
echo "GET NYC FOOD"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 5
echo "completed trace script"
sleep 10
echo "GET NYC FOOD AGAIN"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 140
echo "GET NYC FOOD ONE MORE TIME"
curl http://localhost:8080/cities/new%20york%20city/food
sleep 10
echo "GET WRONG URL LINKS"
curl http://localhost:8080/cities/new%20york%20city/cafe
sleep 44
curl http://localhost:8080/cities/new%20york%20city/build
sleep 94
curl http://localhost:8080/cities/new%20york%20city/
echo
echo
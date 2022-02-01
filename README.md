
# GOLANG NOMINATIM API

Small api in Go that retrieves data from  the [Nominatim API](https://nominatim.org/release-docs/latest/api/Overview/)

# USAGE
- Before using the API you need to have installed docker and docker-compose installed.
- To run the API you need to run the following command:
> docker-compose up --build
- To query about a location you can do a curl, ex:
> curl localhost:9090/api\?q=new%20york
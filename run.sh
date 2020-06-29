./setup-mongo-rs.sh

# running the other services, generating data by default
GENERATE_DATA=true GENERATE_DATA_PERIOD=100 docker-compose up -d
#!/bin/bash

# Prompt user for config file path if not provided as an argument
if [ -z "$1" ]; then
  read -p "Please enter the configuration file path (or press Enter to use the default): " CONFIG_FILE
  CONFIG_FILE=${CONFIG_FILE:-"/Users/admin/Desktop/synnergy_network_demo/synnergy_network_version_1.0/cmd/configs/standard_configs/node_config.yaml"}
else
  CONFIG_FILE=$1
fi

# Function to extract information from YAML config
get_value_from_config() {
  local key=$1
  grep "$key" $CONFIG_FILE | awk '{print $2}'
}

# Function to discover geolocation
discover_geolocation() {
  # Using an external geolocation service to determine latitude and longitude
  GEO_INFO=$(curl -s "https://ipinfo.io/geo")
  LATITUDE=$(echo $GEO_INFO | grep -o '"loc": "[^"]*' | grep -o '[^"]*$' | cut -d',' -f1)
  LONGITUDE=$(echo $GEO_INFO | grep -o '"loc": "[^"]*' | grep -o '[^"]*$' | cut -d',' -f2)

  echo "Discovered Latitude: $LATITUDE, Longitude: $LONGITUDE"
}

# Get node details from the config file
NODE_IP=$(get_value_from_config "node: ip")
NODE_TYPE=$(get_value_from_config "node: type")
NODE_ID=$(get_value_from_config "node: id")
NODE_NAME=$(get_value_from_config "node: name")

# Get geolocation settings from the config file
GEOLOCATION_ENABLED=$(get_value_from_config "geolocation: enabled")
NODE_LATITUDE=$(get_value_from_config "node: latitude")
NODE_LONGITUDE=$(get_value_from_config "node: longitude")

# Function to register the node's location
register_node_location() {
  local node_ip=$1
  local latitude=$2
  local longitude=$3

  echo "Registering location for Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/register_node_location" \
    -H "Content-Type: application/json" \
    -d '{
          "node_id": "'"$NODE_ID"'",
          "latitude": '"$latitude"',
          "longitude": '"$longitude"'
        }'
}

# Function to find the closest node to the current node
find_closest_node() {
  local node_ip=$1
  local latitude=$2
  local longitude=$3

  echo "Finding the closest node to Node: $NODE_NAME ($NODE_ID)"
  curl -X POST "$node_ip/api/network/find_closest_node" \
    -H "Content-Type: application/json" \
    -d '{
          "latitude": '"$latitude"',
          "longitude": '"$longitude"'
        }'
}

# Main geolocation setup process
setup_geolocation() {
  local node_ip=$1
  local latitude=$2
  local longitude=$3

  if [ "$GEOLOCATION_ENABLED" == "true" ]; then
    # If latitude and longitude are not provided, discover them
    if [ -z "$latitude" ] || [ -z "$longitude" ]; then
      echo "Latitude or Longitude not provided. Discovering geolocation..."
      discover_geolocation
      latitude=$LATITUDE
      longitude=$LONGITUDE
    fi

    # Register the node's geolocation
    register_node_location "$node_ip" "$latitude" "$longitude"

    # Find the closest node to this node
    find_closest_node "$node_ip" "$latitude" "$longitude"

    echo "Geolocation setup completed for Node: $NODE_NAME ($NODE_TYPE, ID: $NODE_ID) at IP: $node_ip!"
  else
    echo "Geolocation is disabled for Node: $NODE_NAME ($NODE_ID)"
  fi
}

# Execute the geolocation setup
setup_geolocation "$NODE_IP" "$NODE_LATITUDE" "$NODE_LONGITUDE"

if [[ -z ${SERVICE_NAME} ]]; then
    echo "SERVICE_NAME env var not set"
    exit 1
fi

psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}
if [[ ! -z ${SERVICE_PASSWORD} ]]; then
    echo "Creating service user"
    psql -c "CREATE USER ${SERVICE_NAME} WITH PASSWORD '${SERVICE_PASSWORD}';"
    echo "Service user created"
fi

if [[ ! -z ${ADMIN_SERVICE_PASSWORD} ]]; then
    echo "Creating admin service user"
    psql -c "CREATE USER admin_${SERVICE_NAME} WITH PASSWORD '${ADMIN_SERVICE_PASSWORD}';"
    echo "Admin service user created"
fi

echo "Running init sql script"
psql -f 'docker-entrypoint-initdb.d/up.sql'
echo "Initialization is complete"

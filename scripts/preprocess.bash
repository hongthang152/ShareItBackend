for filename in /k8s/*.yml; do
    envsubst < "${filename}"
done
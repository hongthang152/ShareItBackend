apt-get update \
&& apt-get -y install gettext-base \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

for filename in ./k8s/*.yml; do
    envsubst < "${filename}"
done
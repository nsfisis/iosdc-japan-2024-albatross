container_name=$1

docker container exec $container_name mkdir -p /data/files/audio

docker container cp .external-assets/EX_33.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_34.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_35.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_36.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_37.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_38.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_39.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_40.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_41.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_42.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_43.wav $container_name:/data/files/audio
docker container cp .external-assets/EX_44.wav $container_name:/data/files/audio

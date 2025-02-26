for f in $(find "$PROVISIONER_DIRECTORY/devbox/functions" -maxdepth 1 -type f -name '*.sh'); do
    source $f
done


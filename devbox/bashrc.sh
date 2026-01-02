for f in $(find "$PROVISIONER_DIRECTORY/devbox/completions" -maxdepth 1 -type f -name '*.sh'); do
    source $f
done


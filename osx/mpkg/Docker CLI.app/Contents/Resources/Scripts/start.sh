ISO=$HOME/.docker/machine/cache/boot2docker.iso
VM=dev
DOCKER_MACHINE=/usr/local/bin/docker-machine

unset DYLD_LIBRARY_PATH
unset LD_LIBRARY_PATH

mkdir -p ~/.docker/machine/cache
if [ ! -f $ISO ]; then
  cp /usr/local/share/boot2docker/boot2docker.iso $ISO
fi

machine=$($DOCKER_MACHINE ls -q | grep "^$VM$")
if [ -z $machine ]; then
  echo "Creating Machine $VM..."
   $DOCKER_MACHINE create -d virtualbox --virtualbox-memory 2048 $VM
else
  echo "Machine $VM already exists."
fi

echo "Starting machine $VM..."
$DOCKER_MACHINE start dev

echo "Setting environment variables for machine $VM..."
eval $($DOCKER_MACHINE env dev --shell=bash)
clear

cat << EOF
                        ##         .
                  ## ## ##        ==
               ## ## ## ## ##    ===
           /"""""""""""""""""\___/ ===
      ~~~ {~~ ~~~~ ~~~ ~~~~ ~~~ ~ /  ===- ~~~
           \______ o           __/
             \    \         __/
              \____\_______/

EOF
echo ""
echo "Your shell is configured to use docker with the VM: $VM"
echo "You can ssh into the VM via 'docker-machine ssh $VM'"
echo ""

bash -c "$SHELL"

# NNN
Generate ACI, OCI or Docker images without depending on centralized registries. Pull a base rootfs securely
from any URL and build your own container images.

### Features
* Provide friendly CLI to build and provision different types of images.
* Uses base images from any URL in internet
* Allow to export only the diff between the base image and whatever is provisioned inside the container.

### Tasks
* If running on linux, check whether the kernel supports overlayfs or not. Use host by default.
* If running on OS X, use xhyve by default. Download CoreOS stable and boot it.
* Generate key pair for sshing into CoreOS
* Download base image
* Unpack base image in lowerdir
* Mount empty directory as upperdir
* Copy provisioner inside overlayfs workspace
* Enter chroot inside overlayfs workspace running provisioner
* Generate image metadata based on image format specified.
* Package upperdir content
* Return container back to user

### CLI
```
nnn -base https://github.com/hooklift/stack/releases/download/v1.0.1/hooklift-stack-v1.0.1.tar.gz \
-aci -o mycontainer provisioner.sh
....
....
output: mycontainer.tar.gz
```

Image formats supported:
* Docker (-docker)
* AppC (-aci)
* OpenContainer images (-oci)

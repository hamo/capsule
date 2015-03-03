package kernel

var initShell string = `#!/bin/sh

export PATH=/bin:/sbin:$PATH

mount -t proc proc /proc
mount -t sysfs sysfs /sys
mount -t tmpfs tmpfs /run

mount -t tmpfs mdev /dev
mkdir /dev/pts
mount -t devpts devpts /dev/pts

echo "/sbin/mdev" > /proc/sys/kernel/hotplug
/sbin/mdev -s

modprobe virtio
modprobe virtio_ring
modprobe virtio_pci
modprobe virtio_console

# FIXME: mount root

exec /capsuled
`

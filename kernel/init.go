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

modprobe fscache
modprobe 9pnet
modprobe 9pnet_virtio
modprobe 9p

# mount sysinit
mkdir /sysinit
mount -t 9p -o trans=virtio sysinit /sysinit

# FIXME: mount root
# FIXME: parse root from /proc/cmdline

exec /sysinit/capsuled
`

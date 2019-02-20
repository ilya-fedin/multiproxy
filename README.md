# Multiproxy
Multiproxy is a set of custom scripts and utlities with [3proxy](https://github.com/z3APA3A/3proxy) as the main utlity for multi-uplink (multiple Internet connections) proxying.

## Installation
Installation is possible in two ways:  
* By installing a snap package
* By installing via GNU Make

### Firewall installation
Multiproxy can work both through iptables and firewalld. You can use iptables, but I recommend using it with firewalld. Example of firewalld installation on several major distributions:

Arch Linux:  
```sh
sudo pacman -S firewalld
sudo systemctl enable --now firewalld
```

Debian/Ubuntu:  
```sh
sudo apt-get install firewalld
```

Fedora:  
```sh
sudo dnf install firewalld
```

CentOS (but firewalld is preinstalled on it):  
```sh
sudo yum install firewalld
sudo systemctl enable --now firewalld
```

### Network manager installation
Besides, for correct work Multiproxy needs to take actions when state of network interfaces changes. To do this, the dhcpcd hooks are included. You can write hooks for your network manager or install dhcpcd. Example of dhcpcd installation on several major distributions:

Arch Linux (but dhcpcd is preinstalled on it):  
```sh
sudo pacman -S dhcpcd
sudo systemctl enable --now dhcpcd
```

Debian/Ubuntu:  
```sh
sudo apt-get install dhcpcd5
```

Fedora:  
```sh
sudo dnf install dhcpcd
```

CentOS:  
```sh
sudo yum install dhcpcd
sudo systemctl enable --now dhcpcd
```

### Multiproxy installation via snapd
it's necessary to install [snapd](https://docs.snapcraft.io/installing-snapd/6735) in order to proceed with that method. Example of installation on several major distributions:

With [yay](https://github.com/Jguer/yay) on Arch Linux:  
```sh
yay -S snapd
sudo systemctl enable --now snapd.socket
```

Debian/Ubuntu:  
```sh
sudo apt-get install snapd
```

Fedora:
```sh
sudo dnf install snapd
```

CentOS does not support snapd yet.

It's also necessary to create a symbolic link to enable support of snap packages with classic confinement (such as this one):  
```sh
sudo ln -s /var/lib/snapd/snap /snap
```

Multiproxy installation via snapd, itself:  
```sh
sudo snap install --classic --dangerous multiproxy_*.snap
```

Snap packages are on the Releases page.

If there are no errors, you have successfully installed Multiproxy.

### Multiproxy installation via GNU Make
It's necessary to install make and [docker-compose](https://docs.docker.com/compose/install/) in order to proceed with that method. Example of installation on several major distributions:

Arch Linux:
```sh
sudo pacman -S make docker docker-compose
sudo systemctl enable --now docker
```

Debian/Ubuntu:
```sh
sudo apt-get install make docker docker-compose
```

CentOS:
```sh
sudo dnf install make docker docker-compose
sudo systemctl enable --now docker
```

But I want to note that in the case of Debian/Ubuntu, Fedora and CentOS, it's better to install Docker from the official Docker repository.

It's also necessary to copy directory with the repository to /opt (directory with the Multiproxy files should be /opt/multiproxy) and execute these commands:  
```sh
sudo make
sudo make install-all
```

If there are no errors, you have successfully installed Multiproxy.

## Configuration
### dhcpcd
Primarily, you need to configure your network manager to work correctly with Multiproxy. Below is an example for dhcpcd.

The first thing to do is to prevent all interfaces from installing default route. Open /etc/dhcpcd.conf and add to the end:  
```
nogateway
```

Next, enable default route installation for the main network interface, from which the Internet should go to everything else:  
```
interface enp1s0
gateway
metric 0
```

Also, this block contains the `metric 0` line, which is needed so that the route of this interface takes priority over routes set by the Multiproxy fallback hook for the correct operation of SimpleFailover.

Once you have finished editing dhcpcd.conf, save and close the file.

### 3proxy
To configure 3proxy, open `3proxy.cfg`, add your settings and save it. You shouldn't replace `@USER@`, `@PASSWORD@` or `@PORT@`, but you can uncomment the authentication settings. And, if you need one login for all your instances, you can replace `@USER@` and `@PASSWORD@` with your data, but leaving it commented and adding a new `users` directive would be better.

### Multiproxy
To configure Multiproxy itself, open the `instances` file. This file has the following syntax:  
```
192.168.1.1	enp1s0	0x1	101	501	27001	user1	1111	reconnect	600
192.168.2.1	enp2s0	0x2	102	502	27002	user2	2222	reconnect	600
```
Each line is a 3proxy instance. The file has 10 options. Each option is separated by a tab (not by space!).  

Options description:
1. Gateway IP address of the uplink which will be used with 3proxy instance.  
2. Network interface name of the uplink which will be used with 3proxy instance.  
3. Mark that will be added to the packets of 3proxy instance. It's will be used to route the traffic of 3proxy instance with the specified uplink. **You can just add +1 for each instance.**  
4. Route table of 3proxy instance's uplink. The default route of uplink will be added to this table, and traffic with the mark of that uplink will use this table. **You can just add +1 for each instance.**  
5. UID of 3proxy instance. The UID is used to match the traffic and the mark of 3proxy instance. **You can just add +1 for each instance.**  
6. Port of 3proxy instance. **You can just add +1 for each instance.**  
7. Login name of 3proxy instance.  
8. Login password of 3proxy instance.  
9. Reconnect method of 3proxy instance's uplink. If you use a Huawei HiLink device, you can enable auto-reconnect of your device. There are two methods: `reconnect` and `reboot`.  
10. Reconnect interval of 3proxy instance's uplink, in seconds.

You can specify only gateway IP address, only network interface name or both.

After Multiproxy configuration, it's necessary to apply the settings.

#### Settings applying when installation is done using snapd
```sh
sudo multiproxy.apply-settings
```

If there are no errors, the settings have been applied successfully.

#### Settings applying when installation is done using GNU Make
```sh
sudo /opt/multiproxy/bin/apply-settings
```

If there are no errors, the settings have been applied successfully.

### SimpleFailover
This script checks the Internet connection, disables access through main network interface when connection is lost, and re-enables it when connection is restored. Configuration is optional.

When the script will remove the default route, a default route with higher metric will be used. Default routes of the Multiproxy uplinks (Internet connections) use a number in the `table` column of the `instances` file as the route metric. Thus, when the Internet connection is lost, the uplink with the lowest number in the `table` column will be used. This script doesn't track Internet connection on the Multiproxy uplinks, therefore, if you need failovering of all your uplinks, you should install a more powerful failover.

The location of the SimpleFailover configuration if you have installed via a snap package: `/var/snap/multiproxy/current/default/simplefailover`, if you have installed via GNU Make: `/opt/mutilproxy/default/simplefailover`.

## Uninstallation
### snapd
```sh
sudo snap remove multiproxy
```

### make
```sh
cd /opt/multiproxy
sudo make uninstall-all
cd ..
sudo rm -r multiproxy
```

# Gnocco  
## a small cache of goodness  

Gnocco is a DNS cache with resolver, it is based on the wonderful DNS library from Miek Gieben
[dns](https://github.com/miekg/dns) using the iterative resolver from [Darvaza](https://darvaza.org/resolver).

Gnocco is in very early stages of development with most of its features not implemented yet.

## Quick start

0. Keep in mind that Gnocco is NOT production state
1. Create an user and group (gnocco user and group are advised but not mandatory)
2. Create a configuration space (directory) for gnocco (ie. /etc/gnocco)
3. Move gnocco.conf and roots files in the configuration space
4. Review and modify gnocco.conf to suit your needs
5. Create a directory to hold logs (ie. /var/log/gnocco)
6. Move the gnocco binary to /usr/bin
7. Run `sudo setcap cap_net_bind_service=+ep /usr/bin/gnocco` in order to enable Gnocco to listen on ports &lt; 1024 (ie. 53)
8. Itegrate Gnocco with your init system
